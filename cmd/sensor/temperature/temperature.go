package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensorPublisher "airport-weather/internal/sensor-publisher"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 || len(os.Args[1]) != 3 {
		fmt.Println("Please provide a valid airport name.")
		return
	}
	var airportCode string = os.Args[1]
	client := mqttClient.GetMqttClient("pressure-sensor")
	c := make(chan os.Signal, 1)
	go sensorPublisher.PublishSensorValue(client, 1, "temperature", airportCode, 20)
	<-c
}
