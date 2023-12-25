package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensorPublisher "airport-weather/internal/sensor-publisher"
	"os"
)

func main() {
	client := mqttClient.GetMqttClient("temperature-sensor")
	c := make(chan os.Signal, 1)
	go sensorPublisher.PublishSensorValue(client, 1, "temperature", "CDG", 20)
	go sensorPublisher.PublishSensorValue(client, 1, "temperature", "RAK", 25)
	<-c
}
