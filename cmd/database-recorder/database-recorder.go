package main

import (
	dataType "airport-weather/internal/data-type"
	db "airport-weather/internal/database"
	mqttClient "airport-weather/internal/mqtt-client"
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

var dbClient *mongo.Client

func main() {
	dbClient = db.GetDbClient()
	c := make(chan os.Signal, 1)
	mqttClient.Subscribe("airport/+/sensor/#", mqttClient.GetMqttClient("database-recorder"), subHandler)
	<-c
}

var subHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	storeData(msg.Payload())
}

func storeData(payload []byte) {
	var sensorDataType dataType.DataType
	err := json.Unmarshal(payload, &sensorDataType)

	if err != nil {
		log.Fatal(err)
	}
	doc := bson.M{
		"sensorId":   sensorDataType.SensorId,
		"airportId":  sensorDataType.AirportId,
		"sensorType": sensorDataType.SensorType,
		"value":      sensorDataType.Value,
		"timestamp":  primitive.NewDateTimeFromTime(sensorDataType.Timestamp),
	}

	result, err := dbClient.Database("airports").Collection("weather").InsertOne(context.Background(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Data Inserted Successfully with the ID being:", result.InsertedID)
}
