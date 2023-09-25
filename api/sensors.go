package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ValueData struct {
	Timestamp *time.Time             `json:"timestamp"`
	Value     float64                `json:"value,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

type WaterLevel struct {
	Level     int        `json:"liters"`
	Timestamp *time.Time `json:"timestamp"`
}

// WaterLevelSensorHandler handles requests to retrieve water level sensors in a specific tank
func WaterLevelSensorHandler(w http.ResponseWriter, r *http.Request) {
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

	// Unmarshal the JSON data into a slice of DeviceData
	var tanks []Tank
	err = json.Unmarshal(body, &tanks)
	if err != nil {
		fmt.Println("Error unmarshaling tanks:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Filter the sensors based on tankID and kind = "WaterLevel" in the meta field
	var waterLevelSensors []SensorData
	for _, tank := range tanks {
		if tank.ID == tankID {
			for _, sensor := range tank.Sensors {
				if sensor.Meta.Kind == "WaterLevel" {
					waterLevelSensors = append(waterLevelSensors, sensor)
				}
			}
			break
		}
	}

	// Marshal the water level sensors into JSON
	response, err := json.Marshal(waterLevelSensors)
	if err != nil {
		fmt.Println("Error marshaling water level sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// GetWaterLevelHandler handles requests to retrieve the amount of water in liters for a specific tank
func GetWaterLevelValueHandler(w http.ResponseWriter, r *http.Request) {
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

	// Find the tank with the specified ID
	var targetTank Tank
	for _, tank := range tanks {
		if tank.ID == tankID {
			targetTank = tank
			break
		}
	}

	if targetTank.ID == "" {
		// Tank with the specified ID not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find the water level sensor value for the specified tank
	var waterLevelValue interface{}
	var timestamp *time.Time
	for _, sensor := range targetTank.Sensors {
		if sensor.Meta.Kind == "WaterLevel" {
			waterLevelValue = sensor.Value
			timestamp = sensor.Time
			break
		}
	}

	if waterLevelValue == nil {
		// Water level sensor not found for the specified tank
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Calculate the amount of water in liters
	tankHeight := targetTank.Meta.Settings.Height
	tankCapacity := targetTank.Meta.Settings.Capacity
	sensorValue := waterLevelValue.(float64)

	liters := 0

	if tankCapacity > 0 && tankHeight > 0 {

		calculatedValue := ((tankHeight - sensorValue) / tankHeight) * tankCapacity
		liters = int(math.Round(calculatedValue))

	}

	response := WaterLevel{
		Level:     liters,
		Timestamp: timestamp,
	}

	responseJSONBytes, err := json.Marshal(response)

	if err != nil {
		fmt.Println("Error marshaling water level:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the amount of water in liters into JSON
	// response, err := json.Marshal(map[string]int{"liters": liters})
	if err != nil {
		fmt.Println("Error marshaling water level value:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(responseJSONBytes)
}

func GetWaterLevelHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Send a GET request to localhost/devices/tankID/sensors
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors", tankID))
	if err != nil {
		fmt.Println("Error retrieving sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	sensorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading sensor response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the sensor data into a slice of SensorData
	var sensors []SensorData
	err = json.Unmarshal(sensorBody, &sensors)
	if err != nil {
		fmt.Println("Error unmarshaling sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the water level sensor based on the sensor kind in the meta field
	var waterLevelSensor SensorData
	for _, sensor := range sensors {
		if sensor.Meta.Kind == "WaterLevel" {
			waterLevelSensor = sensor
			break
		}
	}

	// Check if a water level sensor was found
	if waterLevelSensor.ID == "" {
		fmt.Println("Water level sensor not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send a GET request to localhost/devices/tankID/sensors/waterlevel/values
	resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/values", tankID, waterLevelSensor.ID))
	if err != nil {
		fmt.Println("Error retrieving water level values:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	valuesBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading values response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the values data into a slice of ValueData
	var values []SensorData
	err = json.Unmarshal(valuesBody, &values)

	if err != nil {
		fmt.Println("Error unmarshaling values:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send a GET request to localhost/devices
	resp2, err := http.Get("http://localhost/devices")
	if err != nil {
		fmt.Println("Error requesting devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp2.Body)
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

	// Find the tank with the specified ID
	var targetTank Tank
	for _, tank := range tanks {
		if tank.ID == tankID {
			targetTank = tank
			break
		}
	}
	// Check if tank information is available
	if targetTank.ID == "" {
		fmt.Println("Tank information not found")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var waterLevelEntries []WaterLevel
	var timestamp *time.Time
	for _, value := range values {
		sensorValue := value.Value

		timestamp = value.Time
		liters := 0
		if targetTank.Meta.Settings.Height > 0 && targetTank.Meta.Settings.Capacity > 0 {
			calculatedValue := ((targetTank.Meta.Settings.Height - sensorValue.(float64)) / targetTank.Meta.Settings.Height) * targetTank.Meta.Settings.Capacity
			liters = int(math.Round(calculatedValue))
		}

		entry := WaterLevel{
			Level:     liters,
			Timestamp: timestamp,
		}
		waterLevelEntries = append(waterLevelEntries, entry)
	}

	responseJSON := struct {
		WaterLevels []WaterLevel `json:"waterLevels"`
	}{
		WaterLevels: waterLevelEntries,
	}

	// Marshal the response object into JSON
	responseJSONBytes, err := json.Marshal(responseJSON)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(responseJSONBytes)
}

// WaterTemperatureSensorHandler handles requests to retrieve water temperature sensors in a specific tank
func WaterTemperatureSensorHandler(w http.ResponseWriter, r *http.Request) {
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

	// Filter the sensors based on tankID and kind = "WaterThermometer" in the meta field
	var waterTemperatureSensors []SensorData
	for _, tank := range tanks {
		if tank.ID == tankID {
			for _, sensor := range tank.Sensors {
				if sensor.Meta.Kind == "WaterThermometer" {
					waterTemperatureSensors = append(waterTemperatureSensors, sensor)
				}
			}
			break
		}
	}

	// Marshal the water temperature sensors into JSON
	response, err := json.Marshal(waterTemperatureSensors)
	if err != nil {
		fmt.Println("Error marshaling water temperature sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// GetWaterTemperatureValueHandler handles requests to retrieve the value of the water temperature sensor for a specific tank
func GetWaterTemperatureValueHandler(w http.ResponseWriter, r *http.Request) {
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
	var waterTemperatureValue interface{}
	for _, tank := range tanks {
		if tank.ID == tankID {
			for _, sensor := range tank.Sensors {
				if sensor.Meta.Kind == "WaterThermometer" {
					waterTemperatureValue = sensor.Value
					break
				}
			}
			break
		}
	}

	// Marshal the water temperature value into JSON
	response, err := json.Marshal(waterTemperatureValue)
	if err != nil {
		fmt.Println("Error marshaling water temperature value:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// GetWaterLevelHistoryHandler handles requests to retrieve water temperature values for a specific tank
func GetWaterTemperatureHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Send a GET request to localhost/devices/tankID/sensors
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors", tankID))
	if err != nil {
		fmt.Println("Error retrieving sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	sensorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading sensor response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the sensor data into a slice of SensorData
	var sensors []SensorData
	err = json.Unmarshal(sensorBody, &sensors)
	if err != nil {
		fmt.Println("Error unmarshaling sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the water temperature sensor based on the sensor kind in the meta field
	var waterTemperatureSensor SensorData
	for _, sensor := range sensors {
		if sensor.Meta.Kind == "WaterThermometer" {
			waterTemperatureSensor = sensor
			break
		}
	}

	// Check if a water temperature sensor was found
	if waterTemperatureSensor.ID == "" {
		fmt.Println("Water temperature sensor not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send a GET request to localhost/devices/tankID/sensors/watertemperature/values
	resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/values", tankID, waterTemperatureSensor.ID))
	if err != nil {
		fmt.Println("Error retrieving water temperature values:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	valuesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading values response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(valuesBody)
}

// WaterQualitySensorHandler handles requests to retrieve water quality sensors in a specific tank
func WaterQualitySensorHandler(w http.ResponseWriter, r *http.Request) {
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

	// Filter the sensors based on tankID and kind = "WaterPollutantSensor" in the meta field
	var waterQualitySensors []SensorData
	for _, tank := range tanks {
		if tank.ID == tankID {
			for _, sensor := range tank.Sensors {
				if sensor.Meta.Kind == "WaterPollutantSensor" {
					waterQualitySensors = append(waterQualitySensors, sensor)
				}
			}
			break
		}
	}

	// Marshal the water quality sensors into JSON
	response, err := json.Marshal(waterQualitySensors)
	if err != nil {
		fmt.Println("Error marshaling water quality sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// GetWaterQualityValueHandler handles requests to retrieve the value of the water quality sensor for a specific tank
func GetWaterQualityValueHandler(w http.ResponseWriter, r *http.Request) {
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

	// Find the water quality sensor value for the specified tank ID
	var waterQualityValue interface{}
	for _, tank := range tanks {
		if tank.ID == tankID {
			for _, sensor := range tank.Sensors {
				if sensor.Meta.Kind == "WaterPollutantSensor" {
					waterQualityValue = sensor.Value
					break
				}
			}
			break
		}
	}

	// Categorize the water quality based on the value ranges
	var waterQuality string
	if waterQualityValue != nil {
		value, ok := waterQualityValue.(float64)
		if !ok {
			fmt.Println("Error converting water quality value to float64")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if value > 0 && value < 300 {
			waterQuality = "Excellent"
		} else if value >= 300 && value < 900 {
			waterQuality = "Good"
		} else if value >= 900 {
			waterQuality = "Poor"
		} else {
			waterQuality = "Unknown"
		}
	} else {
		waterQuality = "No data available"
	}

	// Create a map for the response
	response := map[string]interface{}{
		"waterQuality": waterQuality,
		"tdsValue":     waterQualityValue,
	}

	// Marshal the response into JSON
	responseBody, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(responseBody)
}

// GetWaterQualityHistoryHandler handles requests to retrieve water quality values for a specific tank
func GetWaterQualityHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Send a GET request to localhost/devices/tankID/sensors
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors", tankID))
	if err != nil {
		fmt.Println("Error retrieving sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	sensorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading sensor response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the sensor data into a slice of SensorData
	var sensors []SensorData
	err = json.Unmarshal(sensorBody, &sensors)
	if err != nil {
		fmt.Println("Error unmarshaling sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the water quality sensor based on the sensor kind in the meta field
	var waterQualitySensor SensorData
	for _, sensor := range sensors {
		if sensor.Meta.Kind == "WaterPollutantSensor" {
			waterQualitySensor = sensor
			break
		}
	}

	// Check if a water quality sensor was found
	if waterQualitySensor.ID == "" {
		fmt.Println("Water quality sensor not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send a GET request to localhost/devices/tankID/sensors/waterquality/values
	resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/values", tankID, waterQualitySensor.ID))
	if err != nil {
		fmt.Println("Error retrieving water quality values:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	valuesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading values response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the values data into a slice of ValueData
	var values []ValueData
	err = json.Unmarshal(valuesBody, &values)
	if err != nil {
		fmt.Println("Error unmarshaling values:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Categorize the water quality values based on the value ranges
	var categorizedValues []map[string]interface{}
	for _, value := range values {
		waterQuality := ""
		v := value.Value

		if v > 0 && v < 300 {
			waterQuality = "Excellent"
		} else if v >= 300 && v < 900 {
			waterQuality = "Good"
		} else if v >= 900 {
			waterQuality = "Poor"
		} else {
			waterQuality = "Unknown"
		}

		categorizedValue := map[string]interface{}{
			"tdsValue":     value.Value,
			"waterQuality": waterQuality,
			"timestamp":    value.Timestamp,
		}

		categorizedValues = append(categorizedValues, categorizedValue)
	}

	// Marshal the categorized values into JSON
	response, err := json.Marshal(categorizedValues)
	if err != nil {
		fmt.Println("Error marshaling categorized values:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

func ChangeWaterLevelAlerts(w http.ResponseWriter, r *http.Request) {
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

	// Find the tank with the specified ID
	var targetTank Tank
	for _, tank := range tanks {
		if tank.ID == tankID {
			targetTank = tank
			break
		}
	}

	if targetTank.ID == "" {
		// Tank with the specified ID not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find the water level sensor value for the specified tank
	// sensorID := ""
	var sensorID string
	for _, sensor := range targetTank.Sensors {
		if sensor.Meta.Kind == "WaterLevel" {
			sensorID = sensor.ID
			break
		}
	}

	// Read the request body
	metaBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send a POST request to localhost/devices/tankID/meta
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/devices/%s/sensors/%s/meta", tankID, sensorID), bytes.NewReader(metaBody))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to text/plain
	req.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	// Check the response status code
	if resp2.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp2.Status)
		// Print the response body for further analysis
		responseBody, _ := ioutil.ReadAll(resp2.Body)
		fmt.Println("Response body:", string(responseBody))
		w.WriteHeader(resp2.StatusCode)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write a success response
	response := map[string]string{
		"message": "Meta field updated successfully",
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(responseBytes)
}

func ChangeWaterTemperatureAlerts(w http.ResponseWriter, r *http.Request) {
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

	// Find the tank with the specified ID
	var targetTank Tank
	for _, tank := range tanks {
		if tank.ID == tankID {
			targetTank = tank
			break
		}
	}

	if targetTank.ID == "" {
		// Tank with the specified ID not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find the water level sensor value for the specified tank
	// sensorID := ""
	var sensorID string
	for _, sensor := range targetTank.Sensors {
		if sensor.Meta.Kind == "WaterThermometer" {
			sensorID = sensor.ID
			break
		}
	}

	// Read the request body
	metaBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send a POST request to localhost/devices/tankID/meta
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/devices/%s/sensors/%s/meta", tankID, sensorID), bytes.NewReader(metaBody))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to text/plain
	req.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	// Check the response status code
	if resp2.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp2.Status)
		// Print the response body for further analysis
		responseBody, _ := ioutil.ReadAll(resp2.Body)
		fmt.Println("Response body:", string(responseBody))
		w.WriteHeader(resp2.StatusCode)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write a success response
	response := map[string]string{
		"message": "Meta field updated successfully",
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(responseBytes)
}

func ChangeWaterQualityAlerts(w http.ResponseWriter, r *http.Request) {
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

	// Find the tank with the specified ID
	var targetTank Tank
	for _, tank := range tanks {
		if tank.ID == tankID {
			targetTank = tank
			break
		}
	}

	if targetTank.ID == "" {
		// Tank with the specified ID not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find the water level sensor value for the specified tank
	// sensorID := ""
	var sensorID string
	for _, sensor := range targetTank.Sensors {
		if sensor.Meta.Kind == "WaterPollutantSensor" {
			sensorID = sensor.ID
			break
		}
	}

	// Read the request body
	metaBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send a POST request to localhost/devices/tankID/meta
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/devices/%s/sensors/%s/meta", tankID, sensorID), bytes.NewReader(metaBody))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to text/plain
	req.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	// Check the response status code
	if resp2.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp2.Status)
		// Print the response body for further analysis
		responseBody, _ := ioutil.ReadAll(resp2.Body)
		fmt.Println("Response body:", string(responseBody))
		w.WriteHeader(resp2.StatusCode)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write a success response
	response := map[string]string{
		"message": "Meta field updated successfully",
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(responseBytes)
}
