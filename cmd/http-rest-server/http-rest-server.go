package main

import (
	_ "airport-weather/cmd/http-rest-server/docs"
	dataType "airport-weather/internal/data-type"
	db "airport-weather/internal/database"
	"context"
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *mongo.Client

func main() {
	log.Println("Starting the server...")
	dbClient = db.GetDbClient()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/{airportIATA}/metric/{metric}", getDataBetweenTwoTimes).Methods("GET")
	r.HandleFunc("/api/v1/{airportIATA}/metric/{metric}/average", getAverageForSingleTypeInDay).Methods("GET")
	r.HandleFunc("/api/v1/{airportIATA}/metrics/average", getAverageForAllTypesInDay).Methods("GET")

	r.HandleFunc("/api/v1/{airportIATA}/metric/{metric}/date-range", getDateInterval).Methods("GET")
	r.HandleFunc("/api/v1/metadata/airports", getAvailableAirportIds).Methods("GET")
	r.HandleFunc("/api/v1/{airportIATA}/available-metrics", getAvailableMetrics).Methods("GET")

	r.PathPrefix("/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	handler := cors.Default().Handler(r)
	log.Println("Server will be listening at", "http://localhost:8080/")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("Error spinning up server", err)
	}
	defer db.CloseDbClient(dbClient)

}

// @Summary Get data between two times
// @Description Get data for a specific metric at an airport between two times
// @Tags Data
// @ID getDataBetweenTwoTimes
// @Param airportIATA path string true "The IATA code of the airport"
// @Param metric path string true "The type of metric (e.g., pressure, temperature, wind-speed)"
// @Param startTime query string true "The start time in RFC3339 format"
// @Param endTime query string true "The end time in RFC3339 format"
// @Produce json
// @Success 200 {array} dataType.DataType "Successful response"
// @Router /api/v1/{airportIATA}/metric/{metric} [get]
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

type SuccessfulAverageResponse struct {
	Average float64 `json:"average"`
	Unit    string  `json:"unit"`
}

// @Summary Get average value of a metric in a day
// @Description Get the average value of a specific metric at an airport for a given date
// @Tags Average
// @ID getAverageForSingleTypeInDay
// @Param airportIATA path string true "The IATA code of the airport"
// @Param metric path string true "The type of metric (e.g., pressure, temperature, wind-speed)"
// @Param date query string true "The date in RFC3339 format"
// @Produce json
// @Success 200 {object} SuccessfulAverageResponse "Successful response"
// @Router /api/v1/{airportIATA}/metric/{metric}/average [get]
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

type AverageAllResponse struct {
	Pressure    *float64 `json:"pressure,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	WindSpeed   *float64 `json:"wind-speed,omitempty"`
}

// @Summary Get average value of all metrics in a day
// @Description Get the average value of all metrics at an airport for a given date
// @Tags Average
// @ID getAverageForAllTypesInDay
// @Param airportIATA path string true "The IATA code of the airport"
// @Param date query string true "The date in RFC3339 format"
// @Produce json
// @Success 200 {object} AverageAllResponse "Successful response"
// @Router /api/v1/{airportIATA}/metrics/average [get]
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

type DateIntervalResponse struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// @Summary Get date interval of a specific metric
// @Description Get the date interval for a specific metric at an airport
// @Tags Metadata
// @ID getDateInterval
// @Param airportIATA path string true "The IATA code of the airport"
// @Param metric path string true "The type of metric (e.g., pressure, temperature, wind-speed)"
// @Produce json
// @Success 200 {object} DateIntervalResponse "Successful response"
// @Router /api/v1/{airportIATA}/metric/{metric}/date-range [get]
func getDateInterval(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	collection := dbClient.Database("airports").Collection("weather")

	filter := bson.M{
		"airportId": strings.ToUpper(params["airportIATA"]),
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var minTime, maxTime time.Time
	var firstIteration = true

	for cursor.Next(context.Background()) {
		var sensorData dataType.DataType
		err := cursor.Decode(&sensorData)
		if err != nil {
			log.Println(err)
		}

		if firstIteration {
			minTime = sensorData.Timestamp
			maxTime = sensorData.Timestamp
			firstIteration = false
		} else {
			if sensorData.Timestamp.Before(minTime) {
				minTime = sensorData.Timestamp
			}

			if sensorData.Timestamp.After(maxTime) {
				maxTime = sensorData.Timestamp
			}
		}
	}

	result := map[string]interface{}{
		"startTime": minTime,
		"endTime":   maxTime,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Error encoding result", http.StatusBadRequest)
	}
}

// @Summary Get all available airport IDs
// @Tags Metadata
// @ID getAvailableAirportIds
// @Produce json
// @Success 200 {array} string "Successful response"
// @Router /api/v1/metadata/airports [get]
func getAvailableAirportIds(w http.ResponseWriter, r *http.Request) {
	collection := dbClient.Database("airports").Collection("weather")
	results, err := collection.Distinct(context.Background(), "airportId", bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Error encoding result", http.StatusBadRequest)
	}
}

// @Summary Get available metrics for a specific airport
// @Tags Metadata
// @ID getAvailableMetrics
// @Param airportIATA path string true "The IATA code of the airport"
// @Produce json
// @Success 200 {array} string "Successful response"
// @Router /api/v1/{airportIATA}/available-metrics [get]
func getAvailableMetrics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	collection := dbClient.Database("airports").Collection("weather")
	filter := bson.M{
		"airportId": strings.ToUpper(params["airportIATA"]),
	}
	results, err := collection.Distinct(context.Background(), "sensorType", filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Error encoding result", http.StatusBadRequest)
	}
}
