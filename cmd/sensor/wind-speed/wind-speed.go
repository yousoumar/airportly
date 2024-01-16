package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensorPublisher "airport-weather/internal/sensor-publisher"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 || len(os.Args[1]) != 3 {
		fmt.Println("Veuillez fournir un nom d'aéroport valide.")
		return
	}
	var airportCode string = os.Args[1]
	client := mqttClient.GetMqttClient("pressure-sensor")
	c := make(chan os.Signal, 1)
	go sensorPublisher.PublishSensorValue(client, 3, "wind-speed", airportCode, 15)
	<-c
}
