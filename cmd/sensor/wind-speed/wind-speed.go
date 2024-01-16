package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensorPublisher "airport-weather/internal/sensor-publisher"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 || len(os.Args[1]) != 3 {
		fmt.Println("Please provide a valid airport IATA code.")
		return
	}
	var airportCode string = os.Args[1]
	topic := fmt.Sprintf("airport/%s/sensor/wind-speed", airportCode)
	client := mqttClient.GetMqttClient("wind-speed-sensor")
	c := make(chan os.Signal, 1)
	go sensorPublisher.PublishSensorValue(topic, client, 3, "wind-speed", airportCode, 15)
	<-c
}
