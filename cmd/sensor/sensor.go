package main

import (
	dataType "airport-weather/internal/data-type"
	mqttClient "airport-weather/internal/mqtt-client"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	sensorId := time.Now().Unix()
	PublishSensorValue(topic, client, sensorId, sensorValueName, airportCode, baseValue)
}

func PublishSensorValue(topic string, client mqtt.Client, sensorId int64, sensorValueName string, airportIata string, baseValue float64) {
	noiseAmplitude := 0.5

	minInterval := 10
	maxInterval := 20
	interval := rand.Intn(maxInterval-minInterval+1) + minInterval

	value := baseValue

	for {
		timestamp := time.Now()
		value += rand.Float64()*2*noiseAmplitude - noiseAmplitude
		sensorDataType := dataType.DataType{SensorId: sensorId, AirportId: airportIata, SensorType: sensorValueName, Value: value, Timestamp: timestamp}
		payload, _ := json.Marshal(sensorDataType)
		mqttClient.Publish(topic, client, payload)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
