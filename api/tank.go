package api

import (
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
