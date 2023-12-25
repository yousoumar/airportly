package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensorPublisher "airport-weather/internal/sensor-publisher"
	"os"
)

func main() {
	client := mqttClient.GetMqttClient("pressure-sensor")
	c := make(chan os.Signal, 1)
	go sensorPublisher.PublishSensorValue(client, 2, "pressure", "CDG", 1000)
	go sensorPublisher.PublishSensorValue(client, 2, "pressure", "RAK", 1005)
	<-c
}
