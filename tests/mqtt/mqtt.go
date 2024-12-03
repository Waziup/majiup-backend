package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	var broker = ""
	var port = 1883

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", broker, port))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)

	// Trap SIGINT (Ctrl+C) and SIGTERM (kill) to gracefully disconnect from MQTT broker
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // Block until a signal is received

	client.Disconnect(250)
	fmt.Println("Disconnected from MQTT broker")
}

func sub(client mqtt.Client) {
	topic := "/#"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
