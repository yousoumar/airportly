package sensor_publisher

import (
	mqttClient "airport-weather/internal/mqtt-client"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func PublishSensorValue(client mqtt.Client, sensorId int, sensorValueName string, airportIata string, baseValue float64) {
	noiseAmplitude := 0.5

	interval := 5
	value := baseValue

	for {
		timestamp := time.Now()
		value += rand.Float64()*2*noiseAmplitude - noiseAmplitude

		mqttClient.Publish(client, "{\"sensorId\":\""+fmt.Sprintf("%d", sensorId)+"\",\"airportId\":\""+airportIata+"\",\"sensorType\":\""+sensorValueName+"\",\"value\":\""+fmt.Sprintf("%.2f", value)+"\",\"timestamp\":\""+timestamp.String()+"\"}")
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
