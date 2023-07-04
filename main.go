package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// TransformedDevice represents a device with its properties
type TransformedDevice struct {
	ID            string         `json:"id" bson:"_id"`
	Name          string         `json:"name" bson:"name"`
	Modified      time.Time      `json:"modified" bson:"modified"`
	Created       time.Time      `json:"created" bson:"created"`
	Notifications []Notification `json:"notifications" bson:"notifications"`
	Location      Location       `json:"location" bson:"location"`
	Geometry      Geometry       `json:"geometry" bson:"geometry"`
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

// validate checks if the geometry values are valid
func (g *Geometry) validate() error {
	// Perform validation checks on the field values
	// Return an error if any validation fails

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

	router.GET("/devices", DeviceProxyHandler)

	// Endpoint to get devices data
	router.GET("/majiup-devices", DeviceHandler)

	// Return devices using a specific ID
	router.GET("/majiup-devices/:deviceID", GetDeviceByIDHandler)

	router.GET("/notifications", ListNotificationsHandler)
	router.GET("/majiup-devices/:deviceID/notifications", ListDeviceNotificationsHandler)
	router.POST("/majiup-devices/:deviceID/notification", CreateNotificationHandler)

	// Endpoint to get all notifications
	// router.GET("/notifications", ListNotificationsHandler)

	// Register the endpoints

	// Endpoint to get notifications for a specific device
	// router.GET("/majiup-devices/:deviceID/notifications", DeviceNotificationsHandler)

	// Endpoint to delete a specific notification by ID
	// router.DELETE("/notification/:notificationID", DeleteNotificationHandler)

	// -----------------------------------------------

	// Endpoint to get all sensors for a specific device
	router.GET("/majiup-devices/:deviceID/sensors", DeviceSensorsHandler)

	// Endpoint to get a specific sensor for a device
	router.GET("/majiup-devices/:deviceID/sensors/:sensor_id", DeviceSensorByIDHandler)

	// Endpoint to create a new sensor for a device
	router.POST("/majiup-devices/:deviceID/sensors", CreateSensorHandler)

	// Endpoint to get the value of a specific sensor for a device
	router.GET("/majiup-devices/:deviceID/sensors/:sensor_id/value", GetSensorValueHandler)

	// Endpoint to get all values of a specific sensor for a device
	router.GET("/majiup-devices/:deviceID/sensors/:sensor_id/values", GetSensorValuesHandler)

	// Endpoint to post a value to a specific sensor for a device
	router.POST("/majiup-devices/:deviceID/sensors/:sensor_id/value", PostSensorValueHandler)

	// Endpoint to delete a specific sensor for a device
	router.DELETE("/majiup-devices/:device_id/sensors/:sensor_id", DeleteSensorHandler)

	// --------------------------------------------------

	// Endpoint to get all actuators for a specific device
	router.GET("/majiup-devices/:deviceID/actuators", DeviceActuatorsHandler)

	// Endpoint to get a specific actuator for a device
	router.GET("/majiup-devices/:deviceID/actuators/:actuator_id", DeviceActuatorByIDHandler)

	// Endpoint to create a new actuator for a device
	router.POST("/majiup-devices/:deviceID/actuators", CreateActuatorHandler)

	// Endpoint to get the value of a specific actuator for a device
	router.GET("/majiup-devices/:deviceID/actuators/:actuator_id/value", GetActuatorValueHandler)

	// Endpoint to get all values of a specific actuator for a device
	router.GET("/majiup-devices/:deviceID/actuators/:actuator_id/values", GetActuatorValuesHandler)

	// Endpoint to post a value to a specific actuator for a device
	router.POST("/majiup-devices/:deviceID/actuators/:actuator_id/value", PostActuatorValueHandler)

	// Endpoint to delete a specific actuator for a device
	router.DELETE("/devices/:device_id/actuators/:actuator_id", DeleteActuatorHandler)

	fmt.Println("Majiup server running at PORT 8080")
	http.ListenAndServe(":8080", router)
}

// NotificationsHandler handles the GET request for /majiup-devices/notifications
func NotificationsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Read the notifications JSON file
	notificationsData, err := ioutil.ReadFile("notifications.json")
	if err != nil {
		fmt.Println("Error reading notifications file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a slice of Notification
	var notifications []Notification
	err = json.Unmarshal(notificationsData, &notifications)
	if err != nil {
		fmt.Println("Error unmarshaling notifications:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the notifications slice into JSON
	response, err := json.Marshal(notifications)
	if err != nil {
		fmt.Println("Error marshaling notifications:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// ListDeviceNotificationsHandler handles requests to list notifications for a specific device
func ListDeviceNotificationsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Get the device ID from the URL parameter
	deviceID := params.ByName("deviceID")

	// Read the devices.json file
	data, err := ioutil.ReadFile("devices.json")
	if err != nil {
		fmt.Println("Error reading devices.json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a slice of Device
	var devices []TransformedDevice
	err = json.Unmarshal(data, &devices)
	if err != nil {
		fmt.Println("Error unmarshaling devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the device with the matching ID
	var targetDevice *TransformedDevice
	for i := range devices {
		if devices[i].ID == deviceID {
			targetDevice = &devices[i]
			break
		}
	}

	// If the device is not found, return an error
	if targetDevice == nil {
		fmt.Println("Device not found:", deviceID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal the device's notifications into JSON
	response, err := json.Marshal(targetDevice.Notifications)
	if err != nil {
		fmt.Println("Error marshaling device notifications:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// CreateNotificationHandler handles requests to create a new notification for a specific device
func CreateNotificationHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Get the device ID from the URL parameter
	deviceID := params.ByName("deviceID")

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into a Notification struct
	var newNotification Notification
	err = json.Unmarshal(body, &newNotification)
	if err != nil {
		fmt.Println("Error unmarshaling request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read the devices.json file
	data, err := ioutil.ReadFile("devices.json")
	if err != nil {
		fmt.Println("Error reading devices.json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a slice of Device
	var devices []TransformedDevice
	err = json.Unmarshal(data, &devices)
	if err != nil {
		fmt.Println("Error unmarshaling devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the device with the matching ID
	var targetDevice *TransformedDevice
	for i := range devices {
		if devices[i].ID == deviceID {
			targetDevice = &devices[i]
			break
		}
	}

	// If the device is not found, return an error
	if targetDevice == nil {
		fmt.Println("Device not found:", deviceID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Append the new notification to the device's notifications
	targetDevice.Notifications = append(targetDevice.Notifications, newNotification)

	// Marshal the updated devices slice into JSON
	updatedData, err := json.Marshal(devices)
	if err != nil {
		fmt.Println("Error marshaling updated devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the updated data to the devices.json file
	err = ioutil.WriteFile("devices.json", updatedData, 0644)
	if err != nil {
		fmt.Println("Error writing to devices.json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeviceProxyHandler handles requests to the /device endpoint
func DeviceProxyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// DeviceHandler handles requests to the /majiup-devices endpoint
// DeviceHandler handles requests to the /majiup-devices endpoint
func DeviceHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	// Unmarshal the JSON data into a slice of TransformedDevice
	var devices []TransformedDevice
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
	transformedDevices := make([]TransformedDevice, len(devices))

	// Transform the devices by extracting the required fields
	for i, device := range devices {
		transformedDevices[i] = TransformedDevice{
			ID:       device.ID,
			Name:     device.Name,
			Modified: device.Modified,
			Created:  device.Created,
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

	err = ioutil.WriteFile("devices.json", body, 0644)
	if err != nil {
		fmt.Println("Error writing devices.json:", err)
	}
}

// GetDeviceByIDHandler handles requests to the /majiup-devices/:deviceID endpoint
func GetDeviceByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")

	// Read the devices.json file
	data, err := ioutil.ReadFile("devices.json")
	if err != nil {
		fmt.Println("Error reading devices.json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a slice of TransformedDevice
	var devices []TransformedDevice
	err = json.Unmarshal(data, &devices)
	if err != nil {
		fmt.Println("Error unmarshaling devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the device with the given device ID
	var device TransformedDevice
	for _, d := range devices {
		if d.ID == deviceID {
			device = d
			break
		}
	}

	// Marshal the device struct into JSON
	response, err := json.Marshal(device)
	if err != nil {
		fmt.Println("Error marshaling device:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// ListNotificationsHandler handles requests to list all notifications
func ListNotificationsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Read the devices.json file
	data, err := ioutil.ReadFile("devices.json")
	if err != nil {
		fmt.Println("Error reading devices.json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a slice of TransformedDevice
	var devices []TransformedDevice
	err = json.Unmarshal(data, &devices)
	if err != nil {
		fmt.Println("Error unmarshaling devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a slice to store all notifications
	var notifications []Notification

	// Iterate over all devices and collect their notifications
	for _, device := range devices {
		notifications = append(notifications, device.Notifications...)
	}

	// Marshal the notifications slice into JSON
	response, err := json.Marshal(notifications)
	if err != nil {
		fmt.Println("Error marshaling notifications:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// DeviceNotificationsHandler handles requests to list notifications for a specific device
func DeviceNotificationsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Read the devices.json file
	data, err := ioutil.ReadFile("devices.json")
	if err != nil {
		fmt.Println("Error reading devices.json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a slice of TransformedDevice
	var devices []TransformedDevice
	err = json.Unmarshal(data, &devices)
	if err != nil {
		fmt.Println("Error unmarshaling devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the device ID from the URL parameter
	deviceID := ps.ByName("deviceID")

	// Find the device with the matching ID
	var device TransformedDevice
	for _, d := range devices {
		if d.ID == deviceID {
			device = d
			break
		}
	}

	// Marshal the device's notifications into JSON
	response, err := json.Marshal(device.Notifications)
	if err != nil {
		fmt.Println("Error marshaling notifications:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(response)
}

// DeviceSensorsHandler handles requests to list all sensors for a specific device
func DeviceSensorsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")

	// Send a GET request to localhost/devices/deviceID/sensors
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors", deviceID))
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// DeviceSensorByIDHandler handles requests to retrieve a specific sensor for a device
func DeviceSensorByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")
	sensorID := ps.ByName("sensor_id")

	// Send a GET request to localhost/devices/deviceID/sensors/sensor_id
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s", deviceID, sensorID))
	if err != nil {
		fmt.Println("Error requesting sensor:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// CreateSensorHandler handles requests to create a new sensor for a device
func CreateSensorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")

	// Send a POST request to localhost/devices/deviceID/sensors
	resp, err := http.Post(fmt.Sprintf("http://localhost/devices/%s/sensors", deviceID), "application/json", r.Body)
	if err != nil {
		fmt.Println("Error creating sensor:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// GetSensorValueHandler handles requests to retrieve the value of a specific sensor for a device
func GetSensorValueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")
	sensorID := ps.ByName("sensor_id")

	// Send a GET request to localhost/devices/deviceID/sensors/sensor_id/value
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/value", deviceID, sensorID))
	if err != nil {
		fmt.Println("Error retrieving sensor value:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// GetSensorValuesHandler handles requests to retrieve all values of a specific sensor for a device
func GetSensorValuesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")
	sensorID := ps.ByName("sensor_id")

	// Send a GET request to localhost/devices/deviceID/sensors/sensor_id/values
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/values", deviceID, sensorID))
	if err != nil {
		fmt.Println("Error retrieving sensor values:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// PostSensorValueHandler handles requests to post a value to a specific sensor for a device
func PostSensorValueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")
	sensorID := ps.ByName("sensor_id")

	// Send a POST request to localhost/devices/deviceID/sensors/sensor_id/value
	resp, err := http.Post(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/value", deviceID, sensorID), "application/json", r.Body)
	if err != nil {
		fmt.Println("Error posting sensor value:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// DeleteSensorHandler handles requests to delete a specific sensor for a device
func DeleteSensorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("device_id")
	sensorID := ps.ByName("sensor_id")

	// Send a DELETE request to localhost/devices/device_id/sensors/sensor_id
	client := http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost/devices/%s/sensors/%s", deviceID, sensorID), nil)
	if err != nil {
		fmt.Println("Error creating delete request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Execute the DELETE request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error deleting sensor:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error deleting sensor:", resp.Status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write([]byte(`{"message": "Sensor deleted successfully"}`))
}

// DeviceActuatorsHandler handles requests to list all actuators for a specific device
func DeviceActuatorsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")

	// Send a GET request to localhost/devices/deviceID/actuators
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/actuators", deviceID))
	if err != nil {
		fmt.Println("Error requesting actuators:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// DeviceActuatorByIDHandler handles requests to retrieve a specific actuator for a device
func DeviceActuatorByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")
	actuatorID := ps.ByName("actuator_id")

	// Send a GET request to localhost/devices/deviceID/actuators/actuator_id
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/actuators/%s", deviceID, actuatorID))
	if err != nil {
		fmt.Println("Error requesting actuator:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// CreateActuatorHandler handles requests to create a new actuator for a device
func CreateActuatorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")

	// Send a POST request to localhost/devices/deviceID/actuators
	resp, err := http.Post(fmt.Sprintf("http://localhost/devices/%s/actuators", deviceID), "application/json", r.Body)
	if err != nil {
		fmt.Println("Error creating actuator:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// GetActuatorValueHandler handles requests to retrieve the value of a specific actuator for a device
func GetActuatorValueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")
	actuatorID := ps.ByName("actuator_id")

	// Send a GET request to localhost/devices/deviceID/actuators/actuator_id/value
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/actuators/%s/value", deviceID, actuatorID))
	if err != nil {
		fmt.Println("Error retrieving actuator value:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// GetActuatorValuesHandler handles requests to retrieve all values of a specific actuator for a device
func GetActuatorValuesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")
	actuatorID := ps.ByName("actuator_id")

	// Send a GET request to localhost/devices/deviceID/actuators/actuator_id/values
	resp, err := http.Get(fmt.Sprintf("http://localhost/devices/%s/actuators/%s/values", deviceID, actuatorID))
	if err != nil {
		fmt.Println("Error retrieving actuator values:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// PostActuatorValueHandler handles requests to post a value to a specific actuator for a device
func PostActuatorValueHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("deviceID")
	actuatorID := ps.ByName("actuator_id")

	// Send a POST request to localhost/devices/deviceID/actuators/actuator_id/value
	resp, err := http.Post(fmt.Sprintf("http://localhost/devices/%s/actuators/%s/value", deviceID, actuatorID), "application/json", r.Body)
	if err != nil {
		fmt.Println("Error posting actuator value:", err)
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

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write(body)
}

// DeleteActuatorHandler handles requests to delete a specific actuator for a device
func DeleteActuatorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deviceID := ps.ByName("device_id")
	actuatorID := ps.ByName("actuator_id")

	// Send a DELETE request to localhost/devices/device_id/actuators/actuator_id
	client := http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost/devices/%s/actuators/%s", deviceID, actuatorID), nil)
	if err != nil {
		fmt.Println("Error creating delete request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Execute the DELETE request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error deleting actuator:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error deleting actuator:", resp.Status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the response writer
	w.Write([]byte(`{"message": "Actuator deleted successfully"}`))
}
