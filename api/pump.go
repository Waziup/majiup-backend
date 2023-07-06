package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
