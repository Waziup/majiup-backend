package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Tank represents a tank with its properties
type Tank struct {
	ID       string       `json:"id" bson:"_id"`
	Name     string       `json:"name" bson:"name"`
	Sensors  []SensorData `json:"sensors"`
	Actuators    []ActuatorData   `json:"actuators"`
	Meta     TankMeta     `json:"meta" bson:"meta"`
	Modified time.Time    `json:"modified" bson:"modified"`
	Created  time.Time    `json:"created" bson:"created"`	
}

type TankMeta struct {
	ReceiveNotification bool         `json:"receivenotifications" bson:"receivenotifications"`
	Notifications       Notification `json:"notifications" bson:"notifications"`
	Location            Location     `json:"location" bson:"location"`
	Settings            Settings     `json:"settings" bson:"settings"`
	Profile				Profile		 `json:"profile" bson:"profile"`
	ActuatorID			string	 	 `json:"actuatorID" bson:"actuatorID"`
	Assigned			bool 		 `json:"assigned" bson:"assigned"`
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
type ActuatorData struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`

	Modified time.Time `json:"modified" bson:"modified"`
	Created  time.Time `json:"created" bson:"created"`

	ActuatorMeta ActuatorMeta `json:"meta" bson:"meta"`

	Time  *time.Time  `json:"time" bson:"time"`
	Value interface{} `json:"value" bson:"value"`
}

type Notification struct {
	Messages []Message `json:"messages" bson:"messages"`
}

type Message struct {
	ID       	int    		`json:"id" bson:"id"`
	TankName 	string 		`json:"tank_name" bson:"tank_name"`
	Date     	string 		`json:"time" bson:"time"`
	Priority 	string 		`json:"priority" bson:"priority"`
	Message  	string 		`json:"message"`
	Read  		bool 		`json:"read_status"`
}

type Location struct {
	Cordinates 	Cordinates `json:"cordinates" bson:"cordinates"`
	Address		string `json:"address" bson:"address"`
}

type Cordinates struct {
	Longitude float64 `json:"longitude" bson:"longitude"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
}

type Profile struct {
	FirstName	string `json:"first_name" bson:"first_name"`
	LastName	string `json:"last_name" bson:"last_name"`
	Username	string `json:"username" bson:"username"`
	Phone 		string `json:"phone" bson:"phone"`
	Email 		string `json:"email" bson:"email"`
	Address		string `json:"address" bson:"address"`
} 

type Settings struct {
	Height   float64 `json:"height" bson:"height"`
	Offset   float64 `json:"offset" bson:"offset"`
	Capacity float64 `json:"capacity" bson:"capacity"`
}

type SensorMeta struct {
	Kind        string  `json:"kind" bson:"kind"`
	Unit        string  `json:"units" bson:"units"`
	CriticalMin float64 `json:"critical_min" bson:"critical_min"`
	CriticalMax float64 `json:"critical_max" bson:"critical_max"`
}

type ActuatorMeta struct {
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

// Get analytics
// This function returns
//  1. Daily average consumption
//  2. Estimated consumption time
//  3. Trend, either UP_TREND or DOWN_TREND in water levels

type Consumption struct {
	Quantity        float64  `json:"quantity" bson:"quantity"`
	Duration        float64  `json:"durarion" bson:"duration"` 
}

type Trend struct {
	Value			float64 `json:"value" bson:"value"`
	AmountConsumed	float64 `json:"amountUsed" bson:"amountUsed"`
	Duration 		float64	`json:"days" bson:"days"`
	Indicator		string	`json:"indicator" bson:"indicator"`
}

type Average struct {
	Hourly 		float64 `json:"hourly" bson:"hourly"`
	Daily 		float64 `json:"daily" bson:"daily"`
	
}

type Analytics struct {
	Average 		Average `json:"average" bson:"average"`
	Trend 			Trend	`json:"trend" bson:"trend"`
	DurationLeft	int	`json:"durationLeft" bson:"durationLeft"`
}

func getConsumption(quantity []WaterLevel ) []Consumption {
	if len(quantity) < 2 {
		return nil
	}

	differences := make([]Consumption, len(quantity)-1)
	
	for i := 0; i < len(quantity)-1; i++ {
		differences[i].Quantity = quantity[i+1].Level - quantity[i].Level
		differences[i].Duration = timeDifference(quantity[i].Timestamp, quantity[i+1].Timestamp)
	}

	return differences
}

func timeDifference(timestamp1 *time.Time, timestamp2 *time.Time) float64 {
	duration := timestamp2.Sub(*timestamp1)

	return duration.Hours()
}

func getMovingAverage(data []WaterLevel, windowSize int) []WaterLevel {
	if windowSize <= 0 {
		return nil
	}
	result := make([]WaterLevel, len(data)-windowSize+1)
	

	for i := 0; i <= len(data)-windowSize; i++ {
		var sum float64
		for j := 0; j < windowSize; j++ {
			sum += data[i+j].Level
		}
		average := sum / float64(windowSize)

		// duration := timeDifference(data[i].Timestamp, data[i+1].Timestamp)
		result[i].Level = average
		if i == 0 {
			result[i].Timestamp = data[i].Timestamp
		} else {
			result[i].Timestamp = data[i+1].Timestamp
		}
	}
	
	return result
}

func getConsumptionAverage(consumption []Consumption, span string) float64 {
	var sumQuantity float64
	var sumDuration float64
	for i := 0; i <= len(consumption) - 1; i++ {
		if consumption[i].Quantity < 0 {
			sumQuantity += math.Abs(consumption[i].Quantity)
		}
		sumDuration += consumption[i].Duration
	}


	avg := sumQuantity / sumDuration

	if span == "hrs" {
		return avg

	} else if span == "mins" {
		return avg / 60
	} else if span == "days" {
		return avg * 24
	}
	return avg
}

func getTrend(consumption []Consumption) Trend {
	var sumQuantity float64
	var amount float64
	var trend Trend
	var duration float64

	for i := 0; i <= len(consumption) - 1; i++ {
		// sumQuantity += consumption[i].Quantity	
		duration += consumption[i].Duration

		if consumption[i].Quantity < 0 {
			amount += math.Abs(consumption[i].Quantity)
		}	
	}

	var countValidForTrend int= 10

	if len(consumption) > countValidForTrend {
		consumption = consumption[len(consumption)-countValidForTrend:]
	} else {
		consumption = consumption
	}

	for i := 0; i <= len(consumption) - 1; i++ {
		sumQuantity += consumption[i].Quantity	
	}

	trend.AmountConsumed = amount
	trend.Value = sumQuantity
	trend.Duration = duration / 24
	if sumQuantity < 0 {
		trend.Indicator = "DOWN"
	} else if sumQuantity > 0 {
		trend.Indicator = "UP"
	} else if sumQuantity == 0 {
		trend.Indicator = "NEUTRAL"
	} 


	return trend

}

func getDurationLeft(consumption []Consumption, currentAmount  float64, ) int {
	avg := getConsumptionAverage(consumption, "hours")

	hoursLeft := int(currentAmount / avg)
	daysLeft := int(hoursLeft / 24)
	
	return daysLeft
}

type FromTo struct {
	From 		string `json:"from" bson:"from"`
	To 			string `json:"to" bson:"to"`
}

func getFromTo(hours time.Duration) FromTo {

	var gap FromTo

	loc, err := time.LoadLocation("Africa/Nairobi") // This location is in EAT timezone (UTC+03:00)

	if err != nil {
		fmt.Println("Error loading location:", err)
		return gap
	}

	
	// Current time in the specified timezone
	now := time.Now().In(loc)
	fmt.Println(now)

	// Time 48 hours ago in the specified timezone
	fortyEightHoursAgo := now.Add(-hours * time.Hour)


	// Format time in the specified format
	from := fortyEightHoursAgo.Format(time.RFC3339)
	to := now.Format(time.RFC3339)

	gap.From = from
	gap.To = to

	return gap
}

func getAnalytics(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	
	from := strings.ReplaceAll(r.URL.Query().Get("from"), " ", "+")
	to := strings.ReplaceAll(r.URL.Query().Get("to"), " ", "+")

	q := u.Query()
	q.Set("from", from)
	q.Set("to", to)
	
	u.RawQuery = q.Encode()

	// Perform the GET request
	resp, err = http.Get(u.String())

	// resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/values", tankID, waterLevelSensor.ID))

	fmt.Println()
	if err != nil {
		fmt.Println("Error retrieving water level values:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		fmt.Println("Unexpected content type:", contentType)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	// fmt.Println(values)
	
	if len(values) < 2 {
		
		// resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/values", tankID, waterLevelSensor.ID))
		// if err != nil {
		// 	fmt.Println("Error retrieving water level values:", err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
	
		// defer resp.Body.Close()
	
		// // Check response status code
		// if resp.StatusCode != http.StatusOK {
		// 	fmt.Println("Unexpected status code:", resp.StatusCode)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
	
		// // Check Content-Type
		// contentType := resp.Header.Get("Content-Type")
		// if contentType != "application/json" {
		// 	fmt.Println("Unexpected content type:", contentType)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
	
		// // Read the response body
		// valuesBody, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Println("Error reading values response body:", err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		
		// // Unmarshal the values data into a slice of ValueData
		// err = json.Unmarshal(valuesBody, &values)
		// if err != nil {
		// 	fmt.Println("Error unmarshaling values:", err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
	}

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
		liters := 0.0
		if targetTank.Meta.Settings.Height > 0 && targetTank.Meta.Settings.Capacity > 0 {
			calculatedValue := ((targetTank.Meta.Settings.Height - (sensorValue.(float64) - targetTank.Meta.Settings.Offset)) / targetTank.Meta.Settings.Height) * targetTank.Meta.Settings.Capacity
			liters = float64(calculatedValue)
		}

		entry := WaterLevel{
			Level:     liters,
			Timestamp: timestamp,
		}
		waterLevelEntries = append(waterLevelEntries, entry)
	}
	
	movingAverage := getMovingAverage(waterLevelEntries, 2)
	consumption := getConsumption(movingAverage)

	averageConsumptionDaily := getConsumptionAverage(consumption,  "days")
	averageConsumptionHourly := getConsumptionAverage(consumption,  "hrs")

	trend := getTrend(consumption)
	durationLeft := getDurationLeft(consumption, waterLevelEntries[len(waterLevelEntries)-1].Level)	

	var analytics Analytics

	if len(consumption) > 2 {
		analytics.Trend = trend
		analytics.Average.Daily = averageConsumptionDaily
		analytics.Average.Hourly = averageConsumptionHourly
		analytics.DurationLeft  =  durationLeft
	}


	// responseJSON := struct {
	// 	WaterLevels []WaterLevel `json:"waterLevels"`
	// }{
	// 	WaterLevels: waterLevelEntries,
	// }

	// Marshal the response object into JSON
	responseJSONBytes, err := json.Marshal(analytics)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[%s] Water analytics: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(responseJSONBytes)

}

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
			Actuators:    tank.Actuators,
			Meta:     tank.Meta,
			Modified: tank.Modified,
			Created:  tank.Created,
		}

		tankHeight := tank.Meta.Settings.Height
		tankCapacity := tank.Meta.Settings.Capacity
		tankOffset := tank.Meta.Settings.Offset

		for _, sensor := range tank.Sensors {

			// Check if the sensor kind is "WaterLevel"
			if sensor.Meta.Kind == "WaterLevel" && tankHeight > 0 && tankCapacity > 0 {
				waterLevelValue := (((tankHeight) - (sensor.Value.(float64)-tankOffset)) / tankHeight) * tankCapacity
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
// func TankLocationPostHandler(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)

// 	tankID := vars["tankID"]

// 	if r.Method == http.MethodGet {
// 		// GET request: Retrieve the location data

// 		// Create a new HTTP client
// 		client := http.Client{}

// 		// Send a GET request to localhost/devices/tankID
// 		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/devices/%s", tankID), nil)
// 		if err != nil {
// 			fmt.Println("Error creating request:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// Set the Accept header to specify JSON response
// 		req.Header.Set("Accept", "application/json")

// 		// Send the request
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println("Error retrieving tank:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		// Read the response body
// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Println("Error reading response body:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// Print the response body for debugging
// 		fmt.Println("Response Body:", string(body))

// 		// Unmarshal the JSON data into a Tank struct
// 		var tank Tank
// 		err = json.Unmarshal(body, &tank)
// 		if err != nil {
// 			fmt.Println("Error unmarshaling tank:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// Get the location data from the meta field
// 		location := tank.Meta.Location

// 		// Marshal the location data into JSON
// 		response, err := json.Marshal(location)
// 		if err != nil {
// 			fmt.Println("Error marshaling location:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// Set the Content-Type header to application/json
// 		w.Header().Set("Content-Type", "application/json")

// 		// Write the JSON response to the response writer
// 		w.Write(response)
// 	} else if r.Method == http.MethodPost {
// 		// POST request: Update the location data

// 		// Read the request body
// 		body, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			fmt.Println("Error reading request body:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// Unmarshal the JSON data into a map
// 		var request map[string]interface{}
// 		err = json.Unmarshal(body, &request)
// 		if err != nil {
// 			fmt.Println("Error unmarshaling request body:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// Extract the latitude and longitude from the request body
// 		latitude, ok := request["latitude"].(float64)
// 		if !ok {
// 			fmt.Println("Invalid latitude value")
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		longitude, ok := request["longitude"].(float64)

// 		address, ok := request["address"]

// 		if !ok {
// 			fmt.Println("Invalid longitude value")
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		// Create a new location object
// 		location := Location{
// 			Cordinates: ,
// 		}

// 		// Create a new tank object with the updated location
// 		tank := Tank{
// 			ID:   tankID,
// 			Meta: TankMeta{Location: location},
// 		}

// 		// Marshal the tank object into JSON
// 		tankData, err := json.Marshal(tank)
// 		if err != nil {
// 			fmt.Println("Error marshaling tank:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// Create a new HTTP client
// 		client := http.Client{}

// 		// Send a POST request to localhost/devices/tankID
// 		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/devices/%s/meta", tankID), bytes.NewBuffer(tankData))
// 		if err != nil {
// 			fmt.Println("Error creating request:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// Set the Content-Type header to specify JSON request body
// 		req.Header.Set("Content-Type", "application/json")

// 		// Send the request
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println("Error updating tank location:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		// Check the response status code
// 		if resp.StatusCode != http.StatusOK {
// 			fmt.Println("Error updating tank location:", resp.Status)
// 			w.WriteHeader(resp.StatusCode)
// 			return
// 		}

// 		// Set the Content-Type header to application/json
// 		w.Header().Set("Content-Type", "application/json")

// 		// Write the success response to the response writer
// 		w.Write([]byte(`{"message": "Location updated successfully"}`))
// 	} else {
// 		// Invalid HTTP method
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }

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
	resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s/values?limit=3", tankID, targetSensor.ID))
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

	log.Printf("[%s] Tank sensor history:", time.Now().Format(time.RFC3339))

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
		"message": "Tank profile updated successfully",
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/devices/%s/meta", tankID), nil)
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
