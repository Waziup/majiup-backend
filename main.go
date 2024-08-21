package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/JosephMusya/majiup-backend/api"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

type SensorMeta struct {
	Kind        string  `json:"kind" bson:"kind"`
	Unit        string  `json:"units" bson:"units"`
	CriticalMin float64 `json:"critical_min" bson:"critical_min"`
	CriticalMax float64 `json:"critical_max" bson:"critical_max"`
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

type PumpMeta struct {
	Kind string `json:"kind" bson:"kind"`
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

type Settings struct {
	Height   float64 `json:"height" bson:"height"`
	Offset   float64 `json:"offset" bson:"offset"`
	Capacity float64 `json:"capacity" bson:"capacity"`
}

type TankMeta struct {
	Settings            Settings     `json:"settings" bson:"settings"`
	Notifications       Notification `json:"notifications" bson:"notifications"`
	Profile				Profile		 `json:"profile" bson:"profile"`

}

type Profile struct {
	FirstName	string `json:"first_name" bson:"first_name"`
	LastName	string `json:"last_name" bson:"last_name"`
	Username	string `json:"username" bson:"username"`
	Phone 		string `json:"phone" bson:"phone"`
	Email 		string `json:"email" bson:"email"`
	Address		string `json:"address" bson:"address"`
} 

// Tank represents a tank with its properties
type Tank struct {
	ID       string       `json:"id" bson:"_id"`
	Name     string       `json:"name" bson:"name"`
	Sensors  []SensorData `json:"sensors"`
	Pumps    []PumpData   `json:"actuators"`
	Modified time.Time    `json:"modified" bson:"modified"`
	Created  time.Time    `json:"created" bson:"created"`
	Meta     TankMeta     `json:"meta" bson:"meta"`
	// Token  	 string    	  `json:"token" bson:"token"`
}

type Gateway struct {
	Token		[]string `json:"token" bson:"token"`
}

func sendMessage(message string, phone string) error {
	// apiKey := os.Getenv("SMS_API_KEY")
	// shortCode := os.Getenv("SMS_SHORTCODE")
	// partnerID := os.Getenv("SMS_PARTNER_ID")

	// fmt.Println("PHONE: ",phone)
	// fmt.Println("APIKEY: ",apiKey)
	// fmt.Println("CODE: ",shortCode)
	// fmt.Println("PARTNER: ",partnerID)

	// Prepare the payload
	payload := map[string]string{
		"apikey":    "5fb19e73763aa97a9fda1a7813dc6a3e",
		"partnerID": "10411",
		"message":   message,
		"shortcode": "TextSMS",
		"mobile":    phone,
	}

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://sms.textsms.co.ke/api/services/sendsms/", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a custom HTTP client with a transport that skips certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func sendPushNotification(expoPushToken string, title string, body string) error {
	message := map[string]interface{}{
		"to":    expoPushToken,
		"sound": "default",
		"title": title,
		"body":  body,
		"data":  map[string]string{"someData": "goes here"},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", "https://exp.host/--/api/v2/push/send", bytes.NewBuffer(messageBytes))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	// Create a custom HTTP client with a transport that skips certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getTokens() []string { // Adjust the return type to []string
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost/device/meta", nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return nil
	}

	// Set the Content-Type and Accept headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error obtaining gateway profile:", err)
		return nil
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error: received non-200 status code %d. Response: %s", resp.StatusCode, string(body))
		return nil
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}

	var gateway Gateway

	err = json.Unmarshal(body, &gateway)
	if err != nil {
		log.Println("Error unmarshaling gateway:", err)
		return nil
	}

	return gateway.Token // Return the list of tokens
}

func getPhone() string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost/device/meta", nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return ""
	}

	// Set the Content-Type and Accept headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error obtaining gateway profile:", err)
		return ""
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error: received non-200 status code %d. Response: %s", resp.StatusCode, string(body))
		return ""
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ""
	}

	var tankMeta TankMeta
	
	err = json.Unmarshal(body, &tankMeta)

	if err != nil {
		log.Println("Error unmarshaling gateway:", err)
		return ""
	}

	return strings.TrimSpace(tankMeta.Profile.Phone)
}
var notified bool = false;

var criticalLevelNotified  bool = false;

var tankFull int = 100;
var tankEmpty int = 20;

var mqttClient mqtt.Client

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

func updateTankMessages(tankID string, newMessage Message) {
	// Send a GET request to get the tank metadata
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/devices/%s/meta", tankID), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var tankMeta TankMeta
	err = json.Unmarshal(body, &tankMeta)
	if err != nil {
		fmt.Println("Error unmarshaling tank metadata:", err)
		return
	}

	// Update the messages in the notifications struct
	// tankMeta.Notifications.Messages = append(tankMeta.Notifications.Messages, newMessage)

	tankMeta.Notifications.Messages = append([]Message{newMessage}, tankMeta.Notifications.Messages...)


	// Marshal the updated struct to JSON
	updatedBody, err := json.Marshal(tankMeta)
	if err != nil {
		fmt.Println("Error marshaling updated tank metadata:", err)
		return
	}

	// Send a POST request to update the tank metadata
	postReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/devices/%s/meta", tankID), bytes.NewBuffer(updatedBody))
	if err != nil {
		fmt.Println("Error creating post request:", err)
		return
	}

	postReq.Header.Set("Content-Type", "application/json")

	postResp, err := client.Do(postReq)
	if err != nil {
		fmt.Println("Error sending post request:", err)
		return
	}
	defer postResp.Body.Close()

	if postResp.StatusCode == http.StatusOK {
		fmt.Println("Tank metadata updated successfully with POST")
	} else {
		fmt.Println("Failed to update tank metadata with POST")
	}
}

func checkValForNotifcation (val float64, tankID string, sensorId string) {

	// Send a GET request to localhost/devices/tankID
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/devices/%s", tankID), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
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
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON data into a Tank struct
	var tank Tank
	err = json.Unmarshal(body, &tank)
	if err != nil {
		fmt.Println("Error unmarshaling tank:", err)
		return
	}

	liters := 0.0

	tankHeight := tank.Meta.Settings.Height
	tankOffset := tank.Meta.Settings.Offset
	tankCapacity := tank.Meta.Settings.Capacity

	calculatedValue := ((tankHeight - (val - tankOffset)) / tankHeight) * tankCapacity
	
	liters = float64(calculatedValue)

	percentage := (liters/tankCapacity)*100

	date := time.Now().Add(3*time.Hour).Format("2006-01-02 15:04:05")

	// Send a GET request to localhost/devices/tankID/sensors
	resp, err = http.Get(fmt.Sprintf("http://localhost/devices/%s/sensors/%s", tankID, sensorId))
	if err != nil {
		fmt.Println("Error retrieving sensors:", err)
		return
	}

	defer resp.Body.Close()

	// Read the response body
	sensorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading sensor response body:", err)
		return
	}

	// Unmarshal the sensor data into a slice of SensorData
	var sensor SensorData
	err = json.Unmarshal(sensorBody, &sensor)
	if err != nil {
		fmt.Println("Error unmarshaling sensors:", err)
		return
	}

	lowerLimit := sensor.Meta.CriticalMin
	upperLimit := sensor.Meta.CriticalMax

	tokens := getTokens()
	phone := getPhone()

	// var difference float64 = percentage - prevValue;
	
	// sending alerts during max and min alerts 
	if percentage <= lowerLimit && !notified  {
		fmt.Println("Sending LOW notification")
		title := fmt.Sprintf("%s is almost empty", tank.Name)
		body := fmt.Sprintf("Water level for %s is at %d%%", tank.Name, int(percentage))
		for _, token := range tokens {
			sendPushNotification(token, title, body)
		}

		message := Message{
			TankName: 	tank.Name,
			Message: 	body,
			Date:    	date,
		}
		updateTankMessages(tank.ID, message)

		// phone := string(0)
		// sendMessage(body, phone)
		notified = true
		return
	} else if percentage >= upperLimit && !notified {
		fmt.Println("Sending HIGH notification")
		title := fmt.Sprintf("%s is almost filled", tank.Name)
		body := fmt.Sprintf("Water level for %s is at %d%%", tank.Name, int(percentage))
		
		for _, token := range tokens {
			sendPushNotification(token, title, body)
		}

		message := Message{
			TankName: tank.Name,
			Message: body,
			Date:    	date,
		}
		updateTankMessages(tank.ID, message)

		// phone := string(0)
		// sendMessage(body, phone)
		notified = true
		return
	} else if (percentage > lowerLimit && percentage < upperLimit ) {
		notified = false
		return
	}
	
	// Send alerts at extreme levels
	if percentage >= float64(tankFull) && !criticalLevelNotified  {
		title := fmt.Sprintf("%s is already full", tank.Name)
		body := fmt.Sprintf("Water level for %s is at %d%%. Turn off the pump.", tank.Name, int(percentage))
		
		for _, token := range tokens {
			sendPushNotification(token, title, body)
		}

		sendMessage(body, phone)

		message := Message{
			TankName: tank.Name,
			Message: body,
			Date:    	date,
		}
		updateTankMessages(tank.ID, message)

		// phone := string(0)
		// sendMessage(body, phone)
		criticalLevelNotified = true
		return
	} else if percentage <= float64(tankEmpty) && !criticalLevelNotified {
		title := fmt.Sprintf("%s is running dry", tank.Name)
		body := fmt.Sprintf("Water level for %s is at %d%%. Turn on the pump", tank.Name, int(percentage))
		
		for _, token := range tokens {
			sendPushNotification(token, title, body)
		}
		
		sendMessage(body, phone)

		message := Message{
			TankName: tank.Name,
			Message: body,
			Date:    	date,
		}
		updateTankMessages(tank.ID, message)

		// phone := string(0)
		// sendMessage(body, phone)
		criticalLevelNotified = true
		return
	} else if (percentage > float64(tankEmpty) && percentage < float64(tankFull) ) {
		criticalLevelNotified = false
		return
	}

}

// initialize checking for device when it goes offline ( a boolean variable for online, true by default)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	// start timeout for 10 minutes after this function executes
	// reset timout before 10 minutes if a message is received (online status to true)
	// if timeout resets before 10 minutes end do nothing
	// if timeout executes upto 10 mins, create a print message offlinev (online status to false)
	// when a message is received, toggle the online status to false
	
	regex := regexp.MustCompile(`^devices/([^/]+)/sensors/([^/]+)/value$`)
	matches := regex.FindStringSubmatch(msg.Topic())

	val := string(msg.Payload())

	floatVal, err := strconv.ParseFloat(val, 64)

	if err != nil {
		fmt.Println("Error parsing float:", err)
		return
	}

	if len(matches) >= 3 {
		deviceID := matches[1]
		sensorID := matches[2]

		checkValForNotifcation( floatVal, deviceID, sensorID)
	}

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Printf("[ MQTT ] Connected\n")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
	connectHandler(client)
}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("[ MQTT ] Subscribed to topic: %s\n", topic)
}

func connectMqtt(broker string, topic string, port int) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", broker, port))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	fmt.Printf("mqtt://%s:%d\n", broker, port)


	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[ MQTT ] failed to connect to MQTT broker: %v", token.Error())
	}

	sub(mqttClient, topic)

	// Trap SIGINT (Ctrl+C) and SIGTERM (kill) to gracefully disconnect from MQTT broker
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals // Block until a signal is received
		mqttClient.Disconnect(250)
		fmt.Println("\n [ MQTT ] Disconnected from MQTT broker")
		os.Exit(0)
	}()

	fmt.Println("[ MQTT ] Connected to MQTT broker")
	return nil
}

func maintainMqttConnection(broker string, topic string, port int) {
	for {
		if mqttClient == nil || !mqttClient.IsConnected() {
			fmt.Println("[ MQTT ] Reconnecting to MQTT broker...")
			err := connectMqtt(broker, topic, port)
			if err != nil {
				fmt.Printf("[ MQTT ] Failed to reconnect to MQTT broker: %v\n", err)
			} else {
				fmt.Println("[ MQTT ] Reconnected to MQTT broker successfully")
			}
		}
		time.Sleep(10 * time.Second)
	}
}


type MqttTopic struct {
	TopicId string `json:"topic" bson:"topic"`
}
	

func getMqttTopics () []MqttTopic {

	resp, err := http.Get("http://localhost/devices")
	if err != nil {
		fmt.Println("Error requesting devices:", err)
		return nil
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// Unmarshal the JSON data into a slice of DeviceData
	var tanks []Tank
	err = json.Unmarshal(body, &tanks)
	if err != nil {
		fmt.Println("Error unmarshaling tanks:", err)
		return nil
	}

	// Filter the sensors based on tankID and kind = "WaterLevel" in the meta field
	var topics []MqttTopic

	for _, tank := range tanks {
		for _, sensor := range tank.Sensors {
			if sensor.Meta.Kind == "WaterLevel" {
				topic := MqttTopic{
					TopicId: fmt.Sprintf("devices/%s/sensors/%s/#", tank.ID, sensor.ID),
				}

				topics = append(topics, topic)
			}
		}
		
	}

	return topics
}

func main() {
	apiRouter := mux.NewRouter()
	api.ApiServe(apiRouter)

	frontendRouter := mux.NewRouter()
	appDir := "serve"

	customHandler := func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path
		ext := strings.ToLower(filepath.Ext(filePath))
		contentType := ""
		switch ext {
		case ".js":
			contentType = "application/javascript"
		case ".css":
			contentType = "text/css"
		}

		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}

		http.ServeFile(w, r, filepath.Join(appDir, filePath))
	}

	frontendRouter.PathPrefix("/assets/").Handler(http.HandlerFunc(customHandler))

	frontendRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlFile := filepath.Join(appDir, "index.html")
		http.ServeFile(w, r, htmlFile)
	})

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/", frontendRouter)
	mainRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	log.Printf("[ SUCCESS ] [ %s ] Majiup running at PORT 8082\n", time.Now().Format(time.RFC3339))

	// Start MQTT connection in a separate goroutine

	topics := getMqttTopics()

	fmt.Println(topics)

	// Start MQTT connection and maintain it in a separate goroutine
	for _, topic := range topics {
		go maintainMqttConnection("localhost", topic.TopicId, 1883)
	}

	// for _, topic := range topics {
	// 	err := connectMqtt("localhost", topic.TopicId, 1883)
	// 	// err := connectMqtt("wazigate.local", topic.TopicId, 1883)
	// 	if err != nil {
	// 		log.Printf("[ MQTT ] Failed to connect to MQTT: %v", err)
	// 		// You can choose to proceed without MQTT or handle this error as per your application's requirements.
	// 	}
	// }

	// Start HTTP server
	err := http.ListenAndServe(":8082", mainRouter)
	if err != nil {
		log.Fatalf("[ HTTP ]Failed to start HTTP server: %v", err)
	}
}
