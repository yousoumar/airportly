package sensor_data_type

import "time"

type DataType struct {
	SensorId   int `json:"sensorId"`
	AirportId  string `json:"airportId"`
	SensorType string `json:"sensorType"`
	Value      float64 `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}