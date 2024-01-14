package sensor_publisher

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensor "airport-weather/internal/sensor-data-type"
	"encoding/json"
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
		sensorDataType := sensor.DataType{SensorId: sensorId, AirportId: airportIata, SensorType: sensorValueName, Value: value, Timestamp: timestamp}
		payload, _ := json.Marshal(sensorDataType)
		mqttClient.Publish("sensor", client, payload)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
