package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"strings"
	"time"
)
import mqttClient "airport-weather/internal/mqtt-client"

type SensorData struct {
	SensorId   string `json:"sensorId"`
	AirportId  string `json:"airportId"`
	SensorType string `json:"sensorType"`
	Value      string `json:"value"`
	Timestamp  string `json:"timestamp"`
}

func main() {
	c := make(chan os.Signal, 1)
	client := mqttClient.GetMqttClient("file-recorder")
	mqttClient.Subscribe(client, subHandler)
	<-c
}

var subHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var data SensorData
	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		log.Fatalln("Error parsing JSON data:", err)
		return
	}

	// Parse and format the timestamp
	timestamp, err := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", strings.Split(data.Timestamp, " m=")[0])
	if err != nil {
		log.Fatalln("Error parsing timestamp:", err)
		return
	}
	formattedDate := timestamp.Format("2006-01-02")
	formattedTime := timestamp.Format("15:04:05")

	// Generate filename
	fileName := fmt.Sprintf("%s-%s.csv", data.AirportId, formattedDate)

	// Open or create the file
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open file", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("Failed to close file", err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write data to CSV
	record := []string{formattedTime, data.SensorType, data.Value}
	if err := writer.Write(record); err != nil {
		log.Fatalln("Error writing to CSV:", err)
	}
}
