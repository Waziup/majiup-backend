package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

type Battery struct {
	Percentage 	interface{} 	`json:"percentage" bson:"percentage"`
	Charging 	bool 			`json:"charging" bson:"charging"`
}

func getCharging(tankID string) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors", tankID))
	if err != nil {
		fmt.Println("Error retrieving sensors:", err)
	}

	defer resp.Body.Close()

	// Read the response body
	sensorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading sensor response body:", err)
	}

	// Unmarshal the sensor data into a slice of SensorData
	var sensors []SensorData
	err = json.Unmarshal(sensorBody, &sensors)
	if err != nil {
		fmt.Println("Error unmarshaling sensors:", err)
	}

	// Find the batt sensor based on the sensor kind in the meta field
	var waterLevelSensor SensorData
	for _, sensor := range sensors {
		if sensor.Meta.Kind == "VoltageSensor" {
			waterLevelSensor = sensor
			break
		}
	}

	// Check if a batt sensor was found
	if waterLevelSensor.ID == "" {
		fmt.Println("Battery sensor not found")
		// var values []SensorData
		// w.WriteHeader(http.StatusNotFound)
		// return
	}	

	baseURL := "http://localhost/devices/%s/sensors/%s/values"
	formattedURL := fmt.Sprintf(baseURL, tankID, waterLevelSensor.ID)

	// Parse the URL
	u, err := url.Parse(formattedURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return false
	}
	q := u.Query()	
	
	u.RawQuery = q.Encode()

	// Perform the GET request
	resp, err = http.Get(u.String())

	if err != nil {
		fmt.Println("Error retrieving batt values:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	valuesBody, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading values response body:", err)
		return false
	}

	// Unmarshal the values data into a slice of ValueData
	var values []SensorData
	err = json.Unmarshal(valuesBody, &values)

	if err != nil {
		return false
	}

	if len(values) >= 2 {
        lastTwo := values[len(values)-2:]
		latest := lastTwo[1].Value
		prev := lastTwo[0].Value

		difference := latest.(float64) - prev.(float64)	
		if difference > 0 {
			return true
		} else {
			return false
		}
	}
	
	return false
}
func getBattInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Send a GET request to localhost/devices
	resp, err := http.Get("http://localhost/devices")
	if err != nil {
		fmt.Println("Error requesting devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a slice of Tank
	var tanks []Tank
	err = json.Unmarshal(body, &tanks)
	if err != nil {
		fmt.Println("Error unmarshaling tanks:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the water temperature sensor value for the specified tank ID
	var battInfo Battery
	for _, tank := range tanks {
		if tank.ID == tankID {
			for _, sensor := range tank.Sensors {
				if sensor.Meta.Kind == "VoltageSensor" {
					battInfo.Percentage = sensor.Value
					break
				}
			}			
			break
		}
	}

	status := getCharging(tankID)
	battInfo.Charging = status

	// Marshal the water temperature value into JSON
	response, err := json.Marshal(battInfo)
	if err != nil {
		fmt.Println("Error marshaling battery info:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// battInfo.Percentage = response

	if err != nil {
		fmt.Println("Error marshaling water temperature value:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[%s] Fetched batt percentage value: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(response)
}