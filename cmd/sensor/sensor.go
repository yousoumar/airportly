package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensorPublisher "airport-weather/internal/sensor-publisher"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) <= 3 {
		fmt.Println("Missing a parameter.")
		return
	}
	if len(os.Args[1]) != 3 {
		fmt.Println("Please provide a valid airport IATA code.")
		return
	}
	if os.Args[2] == "" || os.Args[3] == "" {
		fmt.Println("Empty parameter.")
		return
	}
	airportCode := os.Args[1]
	sensorValueName := os.Args[2]
	baseValue, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Println("Error during parsing the based value")
		return
	}
	topic := fmt.Sprintf("airport/%s/sensor/%s", airportCode, sensorValueName)
	clientName := fmt.Sprintf("sensor-%s-%s", airportCode, sensorValueName)
	client := mqttClient.GetMqttClient(clientName)
	c := make(chan os.Signal, 1)
	sensorPublisher.PublishSensorValue(topic, client, sensorValueName, airportCode, baseValue)
	<-c
}
