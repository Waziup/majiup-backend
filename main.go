package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
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
	Notifications []Notification `json:"notifications" bson:"notifications"`
	Location      Location       `json:"location" bson:"location"`
	Geometry      Geometry       `json:"geometry" bson:"geometry"`
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

// Notification represents a notification with its properties
type Notification struct {
	ID         string `json:"id" bson:"id"`
	Message    string `json:"message"`
	ReadStatus bool   `json:"read_status"`
}

type Location struct {
	Longitude float64 `json:"longitude" bson:"longitude"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
}

type Geometry struct {
	Length   float64 `json:"length" bson:"length"`
	Width    float64 `json:"width" bson:"width"`
	Height   float64 `json:"height" bson:"height"`
	Radius   float64 `json:"radius" bson:"radius"`
	Capacity float64 `json:"capacity" bson:"capacity"`
	Type     string  `json:"type,omitempty" bson:"type"`
}

type SensorMeta struct {
	Kind string `json:"kind" bson:"kind"`
}

type PumpMeta struct {
	Kind string `json:"kind" bson:"kind"`
}

type ValueData struct {
	Timestamp time.Time              `json:"timestamp"`
	Value     int                    `json:"value,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

// validate checks if the geometry values are valid
func (g *Geometry) validate() error {
	// Example validation: Ensure non-negative values for length, width, and height
	if g.Length < 0 || g.Width < 0 || g.Height < 0 {
		return errors.New("geometry dimensions must be non-negative")
	}

	// Example validation: Ensure non-negative value for capacity
	if g.Capacity < 0 {
		return errors.New("geometry capacity must be non-negative")
	}

	// Additional validation checks...

	return nil
}

func main() {

	router := httprouter.New()

	/*----------------------------------TANK ENDPOINTS-------------------------------*/

	// Endpoint to get tanks under majiup
	router.GET("/tanks", TankHandler)

	// Return devices using a specific ID
	router.GET("/tanks/:tankID", GetTankByIDHandler)

	// Endpoint to get all sensors for a specific tank
	router.GET("/tanks/:tankID/tank-sensors", TankSensorHandler)

	/*--------------------------------TANK META ENDPOINTS-------------------------------*/

	/*-----------------------------WATER LEVEL SENSOR ENDPOINTS--------------------------------*/

	// Endpoint to get the water level sensor data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/waterlevel", WaterLevelSensorHandler)

	// Endpoint to get the water level value
	router.GET("/tanks/:tankID/tank-sensors/waterlevel/value", GetWaterLevelValueHandler)

	// Endpoint to get the water level history values
	router.GET("/tanks/:tankID/tank-sensors/waterlevel/values", GetWaterLevelHistoryHandler)

	/*-----------------------------WATER TEMPERATURE SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water temperature sensor data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-temperature", WaterTemperatureSensorHandler)

	// Endpoint to get the water temperature value from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-temperature/value", GetWaterTemperatureValueHandler)

	// Endpoint to get the water temperature history values data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-temperature/values", GetWaterTemperatureHistoryHandler)

	/*-----------------------------WATER QUALITY SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water quality sensor data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-quality", WaterQualitySensorHandler)

	// Endpoint to get the water quality sensor data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-quality/value", GetWaterQualityValueHandler)

	// Endpoint to get the water quality history values from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-quality/values", GetWaterQualityHistoryHandler)

	/*---------------------------------PUMP ENDPOINTS-------------------------------------------*/
	// Endpoint to get the pumps available in the tank (1)
	router.GET("/tanks/:tankID/pumps", TankPumpHandler)

	// Endpoint to get the pump state of the selected tankID
	router.GET("/tanks/:tankID/pumps/state", TankStateHandler)

	// Endpoint to get the pump states of the selected tankID, the length of the array will give the actuation
	router.GET("/tanks/:tankID/pumps/states", TankStateHistoryHandler)

	// Endpoint to get
	router.POST("/tanks/:tankID/pumps/state", TankStatePostHandler)

	fmt.Println("Majiup server running at PORT 8080")
	http.ListenAndServe(":8080", router)
}

// DeviceHandler handles requests to the /tanks endpoint
func TankHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	if len(devices) > 0 {
		devices = devices[1:]
	}

	// Create a new slice to store the transformed devices
	transformedDevices := make([]Tank, len(devices))

	// Transform the devices by extracting the required fields
	for i, tank := range devices {
		transformedDevices[i] = Tank{
			ID:       tank.ID,
			Name:     tank.Name,
			Sensors:  tank.Sensors,
			Pumps:    tank.Pumps,
			Modified: tank.Modified,
			Created:  tank.Created,
		}
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

// TankSensorsHandler handles requests to list all sensors for a specific tank
func TankSensorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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

	// Write the JSON response to the response writer
	w.Write(responseBody)
}

// GetTankByIDHandler handles requests to the /tanks/:tankID endpoint
func GetTankByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

	// Read the devices.json file
	data, err := ioutil.ReadFile("devices.json")
	if err != nil {
		fmt.Println("Error reading devices.json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a slice of Tank
	var tanks []Tank
	err = json.Unmarshal(data, &tanks)
	if err != nil {
		fmt.Println("Error unmarshaling tanks:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the tank with the given tank ID
	var tank Tank
	for _, d := range tanks {
		if d.ID == tankID {
			tank = d
			break
		}
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

	// Write the JSON response to the response writer
	w.Write(response)
}

// WaterLevelSensorHandler handles requests to retrieve water level sensors in a specific tank
func WaterLevelSensorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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

// GetWaterLevelHandler handles requests to retrieve the value of the water level sensor for a specific tank
func GetWaterLevelValueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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

	// Find the water level sensor value for the specified tank ID
	var waterLevelValue interface{}
	for _, tank := range tanks {
		if tank.ID == tankID {
			for _, sensor := range tank.Sensors {
				if sensor.Meta.Kind == "WaterLevel" {
					waterLevelValue = sensor.Value
					break
				}
			}
			break
		}
	}

	// Marshal the water level value into JSON
	response, err := json.Marshal(waterLevelValue)
	if err != nil {
		fmt.Println("Error marshaling water level value:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// GetWaterLevelHistoryHandler handles requests to retrieve water level values for a specific tank
func GetWaterLevelHistoryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(valuesBody)
}

// WaterTemperatureSensorHandler handles requests to retrieve water temperature sensors in a specific tank
func WaterTemperatureSensorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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
func GetWaterTemperatureValueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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
func GetWaterTemperatureHistoryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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
func WaterQualitySensorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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
func GetWaterQualityValueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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
func GetWaterQualityHistoryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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

// WaterQualitySensorHandler handles requests to retrieve water quality sensors in a specific tank
func TankPumpHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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
	var targetTank *Tank
	for _, tank := range tanks {
		if tank.ID == tankID {
			targetTank = &tank
			break
		}
	}

	if targetTank == nil {
		fmt.Println("Tank not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Filter the pumps based on the meta field Kind = "Motor"
	var motorActuators []PumpData
	for _, actuator := range targetTank.Pumps {
		if actuator.PumpMeta.Kind == "Motor" {
			motorActuators = append(motorActuators, actuator)
		}
	}

	// Marshal the motor pumps into JSON
	response, err := json.Marshal(motorActuators)
	if err != nil {
		fmt.Println("Error marshaling motor pumps:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// TankStateHandler handles requests to retrieve the state of an actuator in a specific tank
func TankStateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

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
	var targetTank *Tank
	for _, tank := range tanks {
		if tank.ID == tankID {
			targetTank = &tank
			break
		}
	}

	if targetTank == nil {
		fmt.Println("Tank not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find the state of the actuator in the target tank
	var pumpState interface{}
	for _, actuator := range targetTank.Pumps {
		if actuator.PumpMeta.Kind == "Motor" { // Assuming the actuator name is "pump"
			pumpState = actuator.Value
			break
		}
	}

	if pumpState == nil {
		fmt.Println("Actuator state not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal the actuator state into JSON
	response, err := json.Marshal(pumpState)
	if err != nil {
		fmt.Println("Error marshaling actuator state:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// TankStateHistoryHandler handles requests to retrieve the history of sensor values stored in the actuator values for a specific tank
func TankStateHistoryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

	// Send a GET request to localhost/devices/tankID/sensors
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/actuators", tankID))
	if err != nil {
		fmt.Println("Error retrieving pumps:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	actuatorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading actuator response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the actuator data into a slice of PumpData
	var pumps []PumpData
	err = json.Unmarshal(actuatorBody, &pumps)
	if err != nil {
		fmt.Println("Error unmarshaling pumps:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the actuator with the specified kind (Motor)
	var targetPump PumpData
	for _, actuator := range pumps {
		if actuator.PumpMeta.Kind == "Motor" {
			targetPump = actuator
			break
		}
	}

	if targetPump.ID == "" {
		fmt.Println("Actuator not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send a GET request to localhost/devices/tankID/pumps/actuatorID/values
	resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/actuators/%s/values", tankID, targetPump.ID))
	if err != nil {
		fmt.Println("Error retrieving actuator values:", err)
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

	// Categorize the sensor values based on the value ranges
	var categorizedValues []map[string]interface{}
	for _, value := range values {
		v := value.Value

		categorizedValue := map[string]interface{}{
			"pumpState": v,
			"timestamp": value.Timestamp,
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

// TankStatePostHandler handles requests to update the state value of an actuator in a specific tank
func TankStatePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tankID := ps.ByName("tankID")

	// Send a GET request to localhost/devices/tankID/actuators
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/actuators", tankID))
	if err != nil {
		fmt.Println("Error retrieving actuators:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	actuatorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading actuator response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the actuator data into a slice of PumpData
	var pumps []PumpData
	err = json.Unmarshal(actuatorBody, &pumps)
	if err != nil {
		fmt.Println("Error unmarshaling pumps:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the actuator with the specified kind (Motor)
	var targetPump PumpData
	for _, actuator := range pumps {
		if actuator.PumpMeta.Kind == "Motor" {
			targetPump = actuator
			break
		}
	}

	if targetPump.ID == "" {
		fmt.Println("Actuator not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

	// Update the value of the target pump actuator
	targetPump.Value = request["value"]

	// Marshal the updated actuator state into JSON
	response, err := json.Marshal(targetPump.Value)
	if err != nil {
		fmt.Println("Error marshaling actuator state:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)

	// Perform the POST request to update the state value of the actuator
	actuatorURL := fmt.Sprintf("http://localhost/devices/%s/actuators/%s/value", tankID, targetPump.ID)
	_, err = http.Post(actuatorURL, "applicationl/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error updating actuator state:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
