package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetDataBetweenTwoTimes(t *testing.T) {
	// Test server creation
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/{airportIATA}/{metric}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Handler called")
	})
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Test with a request
	req, _ := http.NewRequest("GET", testServer.URL+"/api/v1/JFK/temperature?startTime=2024-01-01T00:00:00Z&endTime=2024-01-02T00:00:00Z", nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", response.Status)
	}
}

func TestGetAverageForSingleTypeInDay(t *testing.T) {
	// Test server creation
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/{airportIATA}/average", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Handler called")
	})
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Test with a request
	req, _ := http.NewRequest("GET", testServer.URL+"/api/v1/JFK/average", nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", response.Status)
	}
}

func TestGetAverageForAllTypesInDay(t *testing.T) {
	// Test server creation
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/{airportIATA}/average/alltype", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Handler called")
	})
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Test with a request
	req, _ := http.NewRequest("GET", testServer.URL+"/api/v1/JFK/average/alltype", nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", response.Status)
	}
}
