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

// handles requests to retrieve actuators in a specific tank
func TankActuatorHandler(w http.ResponseWriter, r *http.Request) {
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

	// Filter the actuators based on the meta field Kind = "Motor"
	var motorActuators []ActuatorData
	for _, actuator := range targetTank.Actuators {
		if actuator.ActuatorMeta.Kind == "Motor" {
			motorActuators = append(motorActuators, actuator)
		}
	}

	// Marshal the motor actuators into JSON
	response, err := json.Marshal(motorActuators)
	if err != nil {
		fmt.Println("Error marshaling motor actuators:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[%s] Retrieved tank actuators: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(response)
}

// TankStateHandler handles requests to retrieve the state of an actuator in a specific tank
func TankStateHandler(w http.ResponseWriter, r *http.Request) {

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
	var actuatorState interface{}
	for _, actuator := range targetTank.Actuators {
		if actuator.ActuatorMeta.Kind == "Motor" { // Assuming the actuator name is "actuator"
			actuatorState = actuator.Value
			break
		}
	}

	if actuatorState == nil {
		fmt.Println("Actuator state not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal the actuator state into JSON
	response, err := json.Marshal(actuatorState)
	if err != nil {
		fmt.Println("Error marshaling actuator state:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Fetched actuator state: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// TankStateHistoryHandler handles requests to retrieve the history of sensor values stored in the actuator values for a specific tank
func TankStateHistoryHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

	// Send a GET request to localhost/devices/tankID/sensors
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

	// Unmarshal the actuator data into a slice of ActuatorData
	var actuators []ActuatorData
	err = json.Unmarshal(actuatorBody, &actuators)
	if err != nil {
		fmt.Println("Error unmarshaling actuators:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the actuator with the specified kind (Motor)
	var targetActuator ActuatorData
	for _, actuator := range actuators {
		if actuator.ActuatorMeta.Kind == "Motor" {
			targetActuator = actuator
			break
		}
	}

	if targetActuator.ID == "" {
		fmt.Println("Actuator not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send a GET request to localhost/devices/tankID/actuators/actuatorID/values
	resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/actuators/%s/values", tankID, targetActuator.ID))
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
			"actuatorState": v,
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

	log.Printf("[%s] Fetched actuator state history: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(response)
}

// TankStatePostHandler handles requests to update the state value of an actuator in a specific tank
func TankStatePostHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tankID := vars["tankID"]

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

	// Unmarshal the actuator data into a slice of ActuatorData
	var actuators []ActuatorData
	err = json.Unmarshal(actuatorBody, &actuators)
	if err != nil {
		fmt.Println("Error unmarshaling actuators:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the actuator with the specified kind (Motor)
	var targetActuator ActuatorData
	for _, actuator := range actuators {
		if actuator.ActuatorMeta.Kind == "Motor" {
			targetActuator = actuator
			break
		}
	}

	if targetActuator.ID == "" {
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

	// The request body should contain just the value, not a key like "state"
	var value interface{}
	err = json.Unmarshal(body, &value)
	if err != nil {
		fmt.Println("Error unmarshaling request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the value of the target actuator actuator
	targetActuator.Value = value

	// Marshal the updated actuator state into JSON
	response, err := json.Marshal(targetActuator.Value)
	if err != nil {
		fmt.Println("Error marshaling actuator state:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[%s] Actuator status changed: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(response)

	// Perform the POST request to update the state value of the actuator
	actuatorURL := fmt.Sprintf("http://localhost/devices/%s/actuators/%s/value", tankID, targetActuator.ID)
	_, err = http.Post(actuatorURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error updating actuator state:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

