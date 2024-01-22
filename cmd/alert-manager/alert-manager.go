package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	dataType "airport-weather/internal/data-type"
	mqttClient "airport-weather/internal/mqtt-client"
)

type Alert struct {
	AirportID  string
	SensorType string
	Value      float64
	Message    string
	Timestamp  time.Time
}

type AlertMessage struct {
	Condition   func(float64) bool
	MessageFunc func(float64) string
}

func main() {
	c := make(chan os.Signal, 1)
	client := mqttClient.GetMqttClient("alert-manager")
	mqttClient.Subscribe("airport/+/sensor/#", client, subHandler)
	<-c
}

var alertMessages = map[string]AlertMessage{
	"temperature": {
		Condition: func(value float64) bool {
			return value <= 16 || value >= 35
		},
		MessageFunc: func(value float64) string {
			if value <= 16 {
				return fmt.Sprintf("Low temperature detected! It's currently %.2f degrees Celsius. It's cold!", value)
			} else {
				return fmt.Sprintf("Hot temperature detected! It's currently %.2f degrees Celsius. Stay cool!", value)
			}
		},
	},
	"pressure": {
		Condition: func(value float64) bool {
			return value >= 1007
		},
		MessageFunc: func(value float64) string {
			return fmt.Sprintf("High pressure detected! The pressure is %.2f hPa.", value)
		},
	},
	"wind-speed": {
		Condition: func(value float64) bool {
			return value >= 35
		},
		MessageFunc: func(value float64) string {
			return fmt.Sprintf("High wind speed detected! The wind speed is %.2f m/s.", value)
		},
	},
}

var subHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var messageData dataType.DataType
	if err := json.Unmarshal(msg.Payload(), &messageData); err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	handleSensorAlert(client, messageData)
}

func handleSensorAlert(client mqtt.Client, data dataType.DataType) {
	alertMessage, exists := alertMessages[data.SensorType]
	if !exists {
		log.Printf("No alert message defined for sensor type: %s\n", data.SensorType)
		return
	}

	if alertMessage.Condition(data.Value) {
		alert := generateAlert(data, alertMessage)
		alert.alert(client)
	}
}

func generateAlert(data dataType.DataType, messageDetails AlertMessage) Alert {
	return Alert{
		AirportID:  data.AirportId,
		SensorType: data.SensorType,
		Value:      data.Value,
		Message:    messageDetails.MessageFunc(data.Value),
		Timestamp:  time.Now(),
	}
}

func (a *Alert) alert(client mqtt.Client) {
	msgPayload, err := json.Marshal(a)
	if err != nil {
		log.Println("Error encoding alert message:", err)
		return
	}

	topic := fmt.Sprintf("airport/%s/alert/%s", a.AirportID, a.SensorType)
	mqttClient.Publish(topic, client, msgPayload)
	log.Printf("Alert published for airport %s - %s: %s\n", a.AirportID, a.SensorType, a.Message)
}
