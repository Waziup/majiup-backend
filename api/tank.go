package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
