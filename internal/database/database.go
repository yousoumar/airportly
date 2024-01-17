package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func GetDbClient() *mongo.Client {

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

	log.Println("Connected to MongoDB successfully !")

	return client
}
