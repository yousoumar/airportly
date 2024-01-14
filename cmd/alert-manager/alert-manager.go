package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	mqttClient "airport-weather/internal/mqtt-client"

	dataType "airport-weather/internal/sensor-data-type"
)

func main() {
	c := make(chan os.Signal, 1)
	client := mqttClient.GetMqttClient("alert-manager")
	mqttClient.Subscribe("sensor", client, subHandler)
	<-c
}

var subHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var messageData dataType.DataType
	err := json.Unmarshal(msg.Payload(), &messageData)
	if err != nil {
		log.Fatal(err)
	}
	switch messageData.SensorType {
	case "temperature":
		if messageData.Value <= 16 {
			alert(client, messageData)
		}
	case "pressure":
		if messageData.Value >= 1007 {
			alert(client, messageData)
		}
	case "wind":
		if messageData.Value >= 35 {
			alert(client, messageData)
		}
	default:
		fmt.Print("Nothing found!")
	}

}

func alert(client mqtt.Client, msg dataType.DataType) {
	msgPayload, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	mqttClient.Publish("alert", client, msgPayload)
}
