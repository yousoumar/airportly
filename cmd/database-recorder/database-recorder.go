package main

import (
	mqttClient "airport-weather/internal/mqtt-client"
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var dbClient *mongo.Client

func main() {
	dbClient = connectDb()
	client := mqttClient.GetMqttClient("data-processor")
	c := make(chan os.Signal, 1)
	mqttClient.Subscribe(client, subHandler)
	<-c
}

var subHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	// Process and store the received data
	storeData(msg.Payload())
}

func storeData(payload []byte) {
	var measurement Measurement
	err := json.Unmarshal(payload, &measurement)

	doc := bson.M{
		"sensor_id":   measurement.SensorId,
		"airport_id":  measurement.AirportId,
		"sensor_type": measurement.SensorType,
		"value":       measurement.Value,
		"timestamp":   measurement.Timestamp,
	}

	result, err := dbClient.Database("airport").Collection("airport").InsertOne(context.Background(), doc)
	if err != nil {
		log.Fatal(err)
	}

	// Afficher l'ID généré pour le message inséré
	fmt.Println("Message inséré avec succès. ID:", result.InsertedID)

}

func connectDb() *mongo.Client {
	uri := "mongodb://localhost:27017"

	credential := options.Credential{Username: "root", Password: "example"}

	// Configuration des options du client MongoDB
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credential)

	// Connecter le client à MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		defer func() {
			// Déconnecter le client lorsque vous avez terminé
			if err := client.Disconnect(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
	}

	// Vérifier la connexion en pingant la base de données
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connexion à MongoDB établie avec succès!")

	return client
}

type Measurement struct {
	SensorId   string `json:"sensorId"`
	AirportId  string `json:"airportId"`
	SensorType string `json:"sensorType"`
	Value      string `json:"value"`
	Timestamp  string `json:"timestamp"`
}
