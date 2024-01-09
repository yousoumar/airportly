package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	sensor "airport-weather/internal/sensor-data-type"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func main() {
	dbClient = connectDb()
	c := make(chan os.Signal, 1)
	mqttClient.Subscribe("sensor", mqttClient.GetMqttClient("database-recorder"), subHandler)
	<-c
}

var subHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	storeData(msg.Payload())
}

func storeData(payload []byte) {
	var sensorDataType sensor.SensorDataType
	err := json.Unmarshal(payload, &sensorDataType)

	doc := bson.M{
		"sensorId":   sensorDataType.SensorId,
		"airportId":  sensorDataType.AirportId,
		"sensorType": sensorDataType.SensorType,
		"value":      sensorDataType.Value,
		"timestamp":  sensorDataType.Timestamp,
	}

	result, err := dbClient.Database("airports").Collection("weather").InsertOne(context.Background(), doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data Inserted Successfully with the ID being:", result.InsertedID)
}

func connectDb() *mongo.Client {

	credential := options.Credential{Username: "root", Password: "example"}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		defer func() {
			if err := client.Disconnect(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB successfully !")

	return client
}
