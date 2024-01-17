package main

import (
	db "airport-weather/internal/database"
	sensor "airport-weather/internal/sensor-data-type"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *mongo.Client

func main() {
	dbClient = db.GetDbClient()
	r := mux.NewRouter()
	r.HandleFunc("/{airportIATA}/{metric}", getDataBetweenTwoTimes).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error spinning up server", err)
	}
}

func getDataBetweenTwoTimes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	startTimeParam := r.URL.Query().Get("startTime")
	endTimeParam := r.URL.Query().Get("endTime")

	if startTimeParam == "" || endTimeParam == "" {
		http.Error(w, "startTime, and endTime should be provided as query paramters", http.StatusBadRequest)
		return
	}

	decodedStartTime, decodedStartTimeError := url.QueryUnescape(startTimeParam)
	if decodedStartTimeError != nil {
		http.Error(w, "Badly encoded startTime", http.StatusBadRequest)
		return
	}

	decodedEndTime, decodedEndTimeErr := url.QueryUnescape(endTimeParam)

	if decodedEndTimeErr != nil {
		http.Error(w, "Badly encoded endTime", http.StatusBadRequest)
		return
	}

	parsedStartTime, parsedStartTimeErr := time.Parse(time.RFC3339, decodedStartTime)

	if parsedStartTimeErr != nil {
		http.Error(w, "Invalid startTime format (It should be an ISO date format)", http.StatusBadRequest)
		return
	}

	parsedEndTime, parsedEndTimeErr := time.Parse(time.RFC3339, decodedEndTime)
	if parsedEndTimeErr != nil {
		http.Error(w, "Invalid endTime format (It should be an ISO date format)", http.StatusBadRequest)
		return
	}

	collection := dbClient.Database("airports").Collection("weather")

	var data []sensor.DataType
	filter := bson.M{
		"sensorType": params["metric"],
		"airportId":  strings.ToUpper(params["airportIATA"]),
		"timestamp": bson.M{
			"$gte": primitive.NewDateTimeFromTime(parsedStartTime),
			"$lte": primitive.NewDateTimeFromTime(parsedEndTime),
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

			http.Error(w, "Error while closing DB connection", http.StatusBadRequest)
		}
	}(cursor, context.Background())

	for cursor.Next(context.Background()) {
		var sensorData sensor.DataType
		err := cursor.Decode(&sensorData)
		if err != nil {
			log.Println(err)
		}
		data = append(data, sensorData)
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "encode for result", http.StatusBadRequest)
	}
}
