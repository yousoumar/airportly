package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensorPublisher "airport-weather/internal/sensor-publisher"
	"os"
)

func main() {
	client := mqttClient.GetMqttClient("wind-speed-sensor")
	c := make(chan os.Signal, 1)
	go sensorPublisher.PublishSensorValue(client, 3, "wind-speed", "CDG", 20)
	go sensorPublisher.PublishSensorValue(client, 3, "wind-speed", "RAK", 15)
	<-c
}
