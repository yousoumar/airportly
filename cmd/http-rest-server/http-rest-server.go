package main

import (
	dataType "airport-weather/internal/data-type"
	db "airport-weather/internal/database"
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
	r.HandleFunc("/api/v1/{airportIATA}/{metric}", getDataBetweenTwoTimes).Methods("GET")
	r.HandleFunc("/api/v1/{airportIATA}/{metric}/average", getAverageForSingleTypeInDay).Methods("GET")
	r.HandleFunc("/api/v1/{airportIATA}/average/alltype", getAverageForAllTypesInDay).Methods("GET")
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

	var data []dataType.DataType
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

			log.Println(w, "Error while closing DB connection", http.StatusBadRequest)
		}
	}(cursor, context.Background())

	for cursor.Next(context.Background()) {
		var sensorData dataType.DataType
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

func getAverageForSingleTypeInDay(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dateParam := r.URL.Query().Get("date")

	if dateParam == "" {
		http.Error(w, "Date should be provided as a query parameter", http.StatusBadRequest)
		return
	}

	decodedDate, decodedDateError := url.QueryUnescape(dateParam)
	if decodedDateError != nil {
		http.Error(w, "Badly encoded date", http.StatusBadRequest)
		return
	}

	parsedDate, parsedDateErr := time.Parse(time.RFC3339, decodedDate)
	if parsedDateErr != nil {
		http.Error(w, "Invalid date format (It should be in ISO format)", http.StatusBadRequest)
		return
	}

	startTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endTime := startTime.Add(24 * time.Hour).Add(-time.Second)

	collection := dbClient.Database("airports").Collection("weather")

	filter := bson.M{
		"sensorType": params["metric"],
		"airportId":  strings.ToUpper(params["airportIATA"]),
		"timestamp": bson.M{
			"$gte": primitive.NewDateTimeFromTime(startTime),
			"$lte": primitive.NewDateTimeFromTime(endTime),
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
			log.Println(w, "Error while closing DB connection", http.StatusBadRequest)
		}
	}(cursor, context.Background())

	var sum float64
	var count int

	for cursor.Next(context.Background()) {
		var sensorData dataType.DataType
		err := cursor.Decode(&sensorData)
		if err != nil {
			log.Println(err)
		}
		sum += sensorData.Value
		count++
	}

	if count > 0 {
		average := sum / float64(count)
		result := map[string]interface{}{
			"average": average,
			"unit":    params["metric"],
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, "Error encoding result", http.StatusBadRequest)
		}
	} else {
		err = json.NewEncoder(w).Encode(map[string]string{"message": "No data for the specified date and type"})
		if err != nil {
			http.Error(w, "Error encoding result", http.StatusBadRequest)
		}
	}
}

func getAverageForAllTypesInDay(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dateParam := r.URL.Query().Get("date")

	if dateParam == "" {
		http.Error(w, "Date should be provided as a query parameter", http.StatusBadRequest)
		return
	}

	decodedDate, decodedDateError := url.QueryUnescape(dateParam)
	if decodedDateError != nil {
		http.Error(w, "Badly encoded date", http.StatusBadRequest)
		return
	}

	parsedDate, parsedDateErr := time.Parse(time.RFC3339, decodedDate)
	if parsedDateErr != nil {
		http.Error(w, "Invalid date format (It should be in ISO format)", http.StatusBadRequest)
		return
	}

	startTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endTime := startTime.Add(24 * time.Hour).Add(-time.Second)

	collection := dbClient.Database("airports").Collection("weather")

	filter := bson.M{
		"airportId": strings.ToUpper(params["airportIATA"]),
		"timestamp": bson.M{
			"$gte": primitive.NewDateTimeFromTime(startTime),
			"$lte": primitive.NewDateTimeFromTime(endTime),
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
			log.Println(w, "Error while closing DB connection", http.StatusBadRequest)
		}
	}(cursor, context.Background())

	var sum1, sum2, sum3 float64
	var count1, count2, count3 int

	for cursor.Next(context.Background()) {
		var sensorData dataType.DataType
		err := cursor.Decode(&sensorData)
		if err != nil {
			log.Println(err)
		}

		switch sensorData.SensorType {
		case "pressure":
			sum1 += sensorData.Value
			count1++
		case "temperature":
			sum2 += sensorData.Value
			count2++
		case "wind-speed":
			sum3 += sensorData.Value
			count3++
		}
	}
	averages := make(map[string]float64)

	if count1 > 0 {
		averages["pressure"] = sum1 / float64(count1)
	}

	if count2 > 0 {
		averages["temperature"] = sum2 / float64(count2)
	}

	if count3 > 0 {
		averages["wind-speed"] = sum3 / float64(count3)
	}

	err = json.NewEncoder(w).Encode(averages)
	if err != nil {
		http.Error(w, "Error encoding result", http.StatusBadRequest)
	}
}
