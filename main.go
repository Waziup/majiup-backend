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
	Token		string `json:"token" bson:"token"`
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
func getToken () string {
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

	var gateway Gateway
	
	err = json.Unmarshal(body, &gateway)

	if err != nil {
		log.Println("Error unmarshaling gateway:", err)
		return ""
	}

	return gateway.Token
}

var notified bool = false;
// var lastNotificationTime time.Time = time.Now();

// var notificationSent bool = false;
// var prevValue float64 = 0.0

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

	token := getToken()

	// var difference float64 = percentage - prevValue;
	
	// fmt.Println("Difference: ",difference)
	fmt.Println("Percentage: ",percentage)
	fmt.Println("Notified: ",notified)
	fmt.Println("EXPO TOKEN: ", token)
	fmt.Println()

	if percentage <= lowerLimit && !notified  {
		fmt.Println("Sending LOW notification")
		title := fmt.Sprintf("%s is almost empty", tank.Name)
		body := fmt.Sprintf("Water level for %s is at %d%%", tank.Name, int(percentage))
		sendPushNotification(token, title, body)
		notified = true
		return
	} else if percentage >= upperLimit && !notified {
		fmt.Println("Sending HIGH notification")
		title := fmt.Sprintf("%s is almost filled", tank.Name)
		body := fmt.Sprintf("Water level for %s is at %d%%", tank.Name, int(percentage))
		sendPushNotification(token, title, body)
		notified = true
		return
	} else if (percentage > lowerLimit && percentage < upperLimit ) {
		notified = false
		return
	}
	
	// prevValue = percentage;		
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	regex := regexp.MustCompile(`^devices/([^/]+)/sensors/([^/]+)/value$`)
	matches := regex.FindStringSubmatch(msg.Topic())

	val := string(msg.Payload())

	floatVal, err := strconv.ParseFloat(val, 64)



	if err != nil {
		fmt.Println("Error parsing float:", err)
		return
	}

	fmt.Println("Received: ", floatVal)

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


	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[ MQTT ] failed to connect to MQTT broker: %v", token.Error())
	}

	sub(client, topic)

	// Trap SIGINT (Ctrl+C) and SIGTERM (kill) to gracefully disconnect from MQTT broker
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals // Block until a signal is received
		client.Disconnect(250)
		fmt.Println("\n [ MQTT ] Disconnected from MQTT broker")
		os.Exit(0)
	}()

	fmt.Println("[ MQTT ] Connected to MQTT broker")
	return nil
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

	for _, topic := range topics {
		err := connectMqtt("localhost", topic.TopicId, 1883)
		// err := connectMqtt("wazigate.local", topic.TopicId, 1883)
		if err != nil {
			log.Printf("[ MQTT ] Failed to connect to MQTT: %v", err)
			// You can choose to proceed without MQTT or handle this error as per your application's requirements.
		}
	}

	// Start HTTP server
	err := http.ListenAndServe(":8082", mainRouter)
	if err != nil {
		log.Fatalf("[ HTTP ]Failed to start HTTP server: %v", err)
	}
}
