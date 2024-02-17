package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather/37.7749/-122.4194", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/weather/{lat}/{lon}", weatherHandler).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"weather":"Clear","temperature":288.11,"feels_like":285.51}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetCoordinates(t *testing.T) {
	vars := map[string]string{"lat": "37.7749", "lon": "-122.4194"}
	lat, lon := getCoordinates(vars)

	if lat != 37.7749 || lon != -122.4194 {
		t.Errorf("getCoordinates returned unexpected values: got lat=%v, lon=%v; want lat=37.7749, lon=-122.4194", lat, lon)
	}

	// Test with invalid coordinates
	vars = map[string]string{"lat": "invalid", "lon": "-122.4194"}
	lat, lon = getCoordinates(vars)

	if lat != 0 || lon != 0 {
		t.Errorf("getCoordinates with invalid input should return lat=0, lon=0: got lat=%v, lon=%v", lat, lon)
	}
}

func TestGetWeatherData(t *testing.T) {
	// Mock the OpenWeather API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"weather":[{"main":"Clear"}],"main":{"temp":288.11,"feels_like":285.51}}`))
	}))
	defer server.Close()

	apiKey = "test_api_key"
	openWeatherAPIURL = server.URL

	lat, lon := 37.7749, -122.4194
	weatherData, err := getWeatherData(lat, lon)

	if err != nil {
		t.Fatalf("getWeatherData returned an error: %v", err)
	}

	if weatherData.Weather[0].Main != "Clear" || weatherData.Main.Temp != 288.11 || weatherData.Main.FeelsLike != 285.51 {
		t.Errorf("getWeatherData returned unexpected data: got %+v", weatherData)
	}
}
