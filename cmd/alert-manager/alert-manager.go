package main

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	mqttClient "airport-weather/internal/mqtt-client"
)

func main() {
	c := make(chan os.Signal, 1)
	client := mqttClient.GetMqttClient("alert-manager")
	mqttClient.Subscribe("sensor", client, subHandler)
	<-c
}

// TODO: We want this function to send alerts to the MQTT server if some thread are reacheds
var subHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
