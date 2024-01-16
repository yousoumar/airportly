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
	airportCode := os.Args[1]
	topic := fmt.Sprintf("airport/%s/sensor/pressure", airportCode)
	client := mqttClient.GetMqttClient("pressure-sensor-" + airportCode)
	c := make(chan os.Signal, 1)
	go sensorPublisher.PublishSensorValue(topic, client, 2, "pressure", airportCode, 1000)
	<-c
}
