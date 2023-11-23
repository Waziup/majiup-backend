package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Tank represents a tank with its properties
type Tank struct {
	ID       string       `json:"id" bson:"_id"`
	Name     string       `json:"name" bson:"name"`
	Sensors  []SensorData `json:"sensors"`
	Pumps    []PumpData   `json:"actuators"`
	Meta     TankMeta     `json:"meta" bson:"meta"`
	Modified time.Time    `json:"modified" bson:"modified"`
	Created  time.Time    `json:"created" bson:"created"`
}

type TankMeta struct {
	ReceiveNotification bool         `json:"receivenotifications" bson:"receivenotifications"`
	Notifications       Notification `json:"notifications" bson:"notifications"`
	Location            Location     `json:"location" bson:"location"`
	Settings            Settings     `json:"settings" bson:"settings"`
}

//Majiup sensor structure
type SensorData struct {
	ID       string      `json:"id" bson:"id"`
	Name     string      `json:"name" bson:"name"`
	Modified time.Time   `json:"modified" bson:"modified"`
	Created  time.Time   `json:"created" bson:"created"`
	Time     *time.Time  `json:"time" bson:"time"`
	Meta     SensorMeta  `json:"meta" bson:"meta"`
	Value    interface{} `json:"value" bson:"value"`
}

//Majiup actuator structure
type PumpData struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`

	Modified time.Time `json:"modified" bson:"modified"`
	Created  time.Time `json:"created" bson:"created"`

	PumpMeta PumpMeta `json:"meta" bson:"meta"`

	Time  *time.Time  `json:"time" bson:"time"`
	Value interface{} `json:"value" bson:"value"`
}

type Notification struct {
	Messages []Message `json:"messages" bson:"messages"`
}

type Message struct {
	ID       int    `json:"id" bson:"id"`
	TankName string `json:"tank_name" bson:"tank_name"`
	Date     string `json:"time" bson:"time"`
	Priority string `json:"priority" bson:"priority"`
	Message  string `json:"message"`
}

type Location struct {
	Longitude float64 `json:"longitude" bson:"longitude"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
}

type Settings struct {
	Height   float64 `json:"height" bson:"height"`
	Capacity float64 `json:"capacity" bson:"capacity"`
}

type SensorMeta struct {
	Kind        string  `json:"kind" bson:"kind"`
	Unit        string  `json:"units" bson:"units"`
	CriticalMin float64 `json:"critical_min" bson:"critical_min"`
	CriticalMax float64 `json:"critical_max" bson:"critical_max"`
}

type PumpMeta struct {
	Kind string `json:"kind" bson:"kind"`
}

// func AskMajiupCopilot(w http.ResponseWriter, r *http.Request) {

// 	hostHeader := r.Host

// 	parts := strings.Split(hostHeader, ":")
// 	ipAddress := parts[0]

// 	const apiKey = "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"
// 	const apiEndpoint = "https://api.openai.com/v1/engines/text-davinci-003/completions" // Update with the appropriate endpoint

// 	query, err := ioutil.ReadAll(r.Body)

// 	reqUrl := fmt.Sprintf("http://%s:8081/api/v1/tanks", ipAddress)

// 	resp1, err := http.Get(reqUrl)

// 	if err != nil {
// 		fmt.Println("Error requesting devices:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp1.Body.Close()
// 	// Read the response body
// 	body, err := ioutil.ReadAll(resp1.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Unmarshal the JSON data into a slice of Tank
// 	var tanks []Tank
// 	err = json.Unmarshal(body, &tanks)
// 	if err != nil {
// 		fmt.Println("Error unmarshaling devices:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	tankJSON, err := json.Marshal(tanks)
// 	if err != nil {
// 		// Handle the error, e.g., log it or return an error response.
// 		fmt.Println("Error marshaling tanks:", err)
// 		return
// 	}

// 	// fmt.Println(tankJSON)

// 	requestData := map[string]interface{}{
// 		"prompt":            string(query) + "\nThese are the tanks available" + string(tankJSON),
// 		"max_tokens":        100, // Customize this according to your needs
// 		"top_p":             1,
// 		"frequency_penalty": 0.6,
// 		"presence_penalty":  0.8,
// 		"temperature":       0.2,
// 	}

// 	jsonData, err := json.Marshal(requestData)
// 	if err != nil {
// 		fmt.Println("Error marshalling JSON:", err)
// 		return
// 	}

// 	client := &http.Client{}

// 	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}

// 	req.Header.Set("Authorization", "Bearer "+apiKey)
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return
// 	}

// 	var response map[string]interface{}
// 	if err := json.Unmarshal(responseBody, &response); err != nil {
// 		fmt.Println("Error decoding JSON response:", err)
// 		return
// 	}

// 	generatedText := response["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)

// 	w.Header().Set("Content-Type", "application/json")

// 	// Marshal the "response" map to JSON
// 	// jsonResponse, err := json.Marshal(generatedText)
// 	jsonResponse := map[string]string{"reply": generatedText}
// 	jsonBytes, err := json.Marshal(jsonResponse)

// 	if err != nil {
// 		fmt.Println("Error marshaling JSON response:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Write the JSON response to the response writer
// 	w.Write(jsonBytes)
// }

// validate checks if the Settings values are valid
// func (g *Settings) validate() error {
// 	// Example validation: Ensure non-negative value for capacity
// 	if g.Capacity < 0 {
// 		return errors.New("Settings capacity must be non-negative")
// 	}

// 	// Additional validation checks...

// 	return nil
// }

// TankHandler handles requests to the /tanks endpoint
func TankHandler(w http.ResponseWriter, r *http.Request) {
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
	var devices []Tank
	err = json.Unmarshal(body, &devices)
	if err != nil {
		fmt.Println("Error unmarshaling devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove the first element from the devices slice
	// if len(devices) > 0 {
	devices = devices[1:]
	// }

	// Create a new slice to store the transformed devices
	transformedDevices := make([]Tank, len(devices))

	for i, tank := range devices {
		// fmt.Println(tank.Sensors)
		var sensorEntry []SensorData

		transformedDevices[i] = Tank{
			ID:       tank.ID,
			Name:     tank.Name,
			Sensors:  tank.Sensors,
			Pumps:    tank.Pumps,
			Meta:     tank.Meta,
			Modified: tank.Modified,
			Created:  tank.Created,
		}

		tankHeight := tank.Meta.Settings.Height
		tankCapacity := tank.Meta.Settings.Capacity

		for _, sensor := range tank.Sensors {

			// Check if the sensor kind is "WaterLevel"
			if sensor.Meta.Kind == "WaterLevel" && tankHeight > 0 && tankCapacity > 0 {
				waterLevelValue := ((tankHeight - sensor.Value.(float64)) / tankHeight) * tankCapacity
				sensor.Value = int(waterLevelValue)
			}

			if sensor.Meta.Kind == "WaterPollutantSensor" {
				var waterQuality string
				// WaterQualityValue := sensor.Value
				if sensor.Value != nil {
					value, ok := sensor.Value.(float64)
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
					sensor.Value = waterQuality
				}
			}

			// fmt.Println(transformedDevices[i])
			sensorEntry = append(sensorEntry, sensor)

			// 	// Update the sensor in the current tank
			// transformedDevices[i].Sensors = sensor
		}

		log.Printf("[%s] Fetched tanks: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

		transformedDevices[i].Sensors = sensorEntry

	}

	// Marshal the transformed devices slice into JSON
	response, err := json.Marshal(transformedDevices)
	if err != nil {
		fmt.Println("Error marshaling devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)

}

// GetTankByIDHandler handles requests to the /tanks/:tankID endpoint
func GetTankByIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Send a GET request to localhost/devices/tankID
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/devices/%s", tankID), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type and Accept headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
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

	// Unmarshal the JSON data into a Tank struct
	var tank Tank
	err = json.Unmarshal(body, &tank)
	if err != nil {
		fmt.Println("Error unmarshaling tank:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the tank struct into JSON
	response, err := json.Marshal(tank)
	if err != nil {
		fmt.Println("Error marshaling tank:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[%s] Fetched tank %s: %s %s", time.Now().Format(time.RFC3339), tankID, r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(response)
}

// TankSensorsHandler handles requests to list all sensors for a specific tank
func TankSensorHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Send a GET request to localhost/devices/tankID/sensors
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors", tankID))
	if err != nil {
		fmt.Println("Error requesting sensors:", err)
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

	// Unmarshal the response body into a slice of SensorData
	var sensors []SensorData
	err = json.Unmarshal(body, &sensors)
	if err != nil {
		fmt.Println("Error unmarshaling sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the SensorData slice into JSON
	responseBody, err := json.Marshal(sensors)
	if err != nil {
		fmt.Println("Error marshaling sensors:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[%s] Fetched tank sensors: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(responseBody)
}

func TankLocationHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Create a new HTTP client
	client := http.Client{}

	// Send a GET request to retrieve the tank data
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/devices/%s", tankID), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Accept header to specify JSON response
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error requesting tank:", err)
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

	// Unmarshal the JSON data into a Tank struct
	tank := Tank{
		// ID:   tankID,
		Meta: TankMeta{Location: Location{}},
	}
	err = json.Unmarshal(body, &tank)
	if err != nil {
		fmt.Println("Error unmarshaling tank:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the location data from the meta field
	location := tank
	fmt.Println(location)
	// fmt.Println(string(body))

	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Marshal the location data into JSON
	response, err := json.Marshal(location)
	if err != nil {
		fmt.Println("Error marshaling location:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// TankLocationHandler handles requests to retrieve or update the location data of a specific tank
func TankLocationPostHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

	if r.Method == http.MethodGet {
		// GET request: Retrieve the location data

		// Create a new HTTP client
		client := http.Client{}

		// Send a GET request to localhost/devices/tankID
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/devices/%s", tankID), nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set the Accept header to specify JSON response
		req.Header.Set("Accept", "application/json")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error retrieving tank:", err)
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

		// Print the response body for debugging
		fmt.Println("Response Body:", string(body))

		// Unmarshal the JSON data into a Tank struct
		var tank Tank
		err = json.Unmarshal(body, &tank)
		if err != nil {
			fmt.Println("Error unmarshaling tank:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the location data from the meta field
		location := tank.Meta.Location

		// Marshal the location data into JSON
		response, err := json.Marshal(location)
		if err != nil {
			fmt.Println("Error marshaling location:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON response to the response writer
		w.Write(response)
	} else if r.Method == http.MethodPost {
		// POST request: Update the location data

		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Unmarshal the JSON data into a map
		var request map[string]interface{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			fmt.Println("Error unmarshaling request body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Extract the latitude and longitude from the request body
		latitude, ok := request["latitude"].(float64)
		if !ok {
			fmt.Println("Invalid latitude value")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		longitude, ok := request["longitude"].(float64)
		if !ok {
			fmt.Println("Invalid longitude value")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create a new location object
		location := Location{
			Latitude:  latitude,
			Longitude: longitude,
		}

		// Create a new tank object with the updated location
		tank := Tank{
			ID:   tankID,
			Meta: TankMeta{Location: location},
		}

		// Marshal the tank object into JSON
		tankData, err := json.Marshal(tank)
		if err != nil {
			fmt.Println("Error marshaling tank:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create a new HTTP client
		client := http.Client{}

		// Send a POST request to localhost/devices/tankID
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/devices/%s/meta", tankID), bytes.NewBuffer(tankData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header to specify JSON request body
		req.Header.Set("Content-Type", "application/json")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error updating tank location:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error updating tank location:", resp.Status)
			w.WriteHeader(resp.StatusCode)
			return
		}

		// Set the Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the success response to the response writer
		w.Write([]byte(`{"message": "Location updated successfully"}`))
	} else {
		// Invalid HTTP method
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type SensorHistory struct {
	WaterLevel       []ValueData `json:"waterLevel"`
	WaterTemperature []ValueData `json:"waterTemperature"`
	WaterQuality     []ValueData `json:"waterQuality"`
}

func GetSensorHistoryHandler(w http.ResponseWriter, r *http.Request) {

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

	// Prepare the sensor history struct
	sensorHistory := SensorHistory{}

	// Find and populate water level history
	waterLevelHistory, err := getSensorHistory(tankID, "WaterLevel")
	if err != nil {
		fmt.Println("Error retrieving water level history:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sensorHistory.WaterLevel = waterLevelHistory

	// Find and populate water temperature history
	waterTemperatureHistory, err := getSensorHistory(tankID, "WaterThermometer")
	if err != nil {
		fmt.Println("Error retrieving water temperature history:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sensorHistory.WaterTemperature = waterTemperatureHistory

	// Find and populate water quality history
	waterQualityHistory, err := getSensorHistory(tankID, "WaterPollutantSensor")
	if err != nil {
		fmt.Println("Error retrieving water quality history:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sensorHistory.WaterQuality = waterQualityHistory

	// Marshal the sensor history into JSON
	response, err := json.Marshal(sensorHistory)
	if err != nil {
		fmt.Println("Error marshaling sensor history:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[%s] Fetched all tank sensor history: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(response)
}

type WaterQualityValueData struct {
	Value        float64   `json:"value"`
	Timestamp    time.Time `json:"timestamp"`
	WaterQuality string    `json:"waterQuality"`
}

func getSensorHistory(tankID, sensorKind string) ([]ValueData, error) {
	// Send a GET request to localhost/devices/tankID/sensors
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors", tankID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	sensorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the sensor data into a slice of SensorData
	var sensors []SensorData
	err = json.Unmarshal(sensorBody, &sensors)
	if err != nil {
		return nil, err
	}

	// Find the sensor based on the sensor kind in the meta field
	var targetSensor SensorData
	for _, sensor := range sensors {
		if sensor.Meta.Kind == sensorKind {
			targetSensor = sensor
			break
		}
	}

	// Check if the target sensor was found
	if targetSensor.ID == "" {
		return nil, fmt.Errorf("%s sensor not found", sensorKind)
	}

	// Send a GET request to localhost/devices/tankID/sensors/sensorID/values
	resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/values", tankID, targetSensor.ID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	valuesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the values data into a slice of ValueData
	var values []ValueData
	err = json.Unmarshal(valuesBody, &values)
	if err != nil {
		return nil, err
	}

	// Assign the correct timestamp to each value
	for i := range values {
		values[i].Timestamp = targetSensor.Time
	}

	log.Printf("[%s] Tank sensor histoey:", time.Now().Format(time.RFC3339))

	return values, nil
}

func ChangeNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Read the new name from the request body
	newTankName, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send a POST request to localhost/devices/tankID/name
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/devices/%s/name", tankID), bytes.NewReader(newTankName))
	if err != nil {
		fmt.Println("Error creating request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to text/plain
	req.Header.Set("Content-Type", "text/plain")

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp.Status)
		w.WriteHeader(resp.StatusCode)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write a success response
	response := map[string]string{
		"message": "Tank name changed successfully",
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Tank name changed: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	w.Write(responseBytes)
}

func postMetaField(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send a POST request to localhost/devices/tankID/meta
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/devices/%s/meta", tankID), bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to text/plain
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp.Status)
		// Print the response body for further analysis
		responseBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response body:", string(responseBody))
		w.WriteHeader(resp.StatusCode)
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

	log.Printf("[%s] Tank meta field updated successfully: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	w.Write(responseBytes)
}

func getMetaFields(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/devices/%s", tankID), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Accept header
	req.Header.Set("Accept", "application/json")

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
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

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp.Status)
		w.WriteHeader(resp.StatusCode)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[%s] Fetched tank meta fields: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the response body to the response writer
	w.Write(body)
}

func DeleteTank(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Create a new DELETE request
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost/devices/%s", tankID), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending DELETE request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp.Status)
		w.WriteHeader(resp.StatusCode)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write a success response
	response := map[string]string{
		"message": "Tank deleted successfully",
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Tank deleted successfuly: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	w.Write(responseBytes)
}
