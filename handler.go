package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var openWeatherAPIURL = "http://api.openweathermap.org/data/2.5/weather"

var apiKey = "c7fae99a3b34958ed8ac00ba29a11ed" // Replace with your OpenWeather API key

type WeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
	} `json:"main"`
	Name string `json:"name"`
	Date int64  `json:"dt"`
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lat, lon := getCoordinates(vars)

	if lat == 0 || lon == 0 {
		http.Error(w, "Invalid coordinates", http.StatusBadRequest)
		return
	}

	weatherData, err := getWeatherData(lat, lon)
	if err != nil {
		http.Error(w, "Unable to fetch weather data", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"weather":     weatherData.Weather[0].Main,
		"temperature": weatherData.Main.Temp,
		"feels_like":  weatherData.Main.FeelsLike,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getCoordinates(vars map[string]string) (float64, float64) {
	lat, lon := vars["lat"], vars["lon"]
	// Convert lat and lon to float64
	latVal, err1 := strconv.ParseFloat(lat, 64)
	lonVal, err2 := strconv.ParseFloat(lon, 64)

	if err1 != nil || err2 != nil {
		return 0, 0
	}
	return latVal, lonVal
}

func getWeatherData(lat, lon float64) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", openWeatherAPIURL, lat, lon, apiKey)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf(response.Status)
	}

	defer response.Body.Close()

	var weatherData WeatherResponse
	if err := json.NewDecoder(response.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}
