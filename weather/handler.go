package weather

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"strconv"
)

var GetCoordinates = getCoordinates
var GetWeatherData = getWeatherData
var openWeatherAPIURL string
var apiKey string
var handlerLog = logrus.New()

//----------------------------------------------------------------------------------------------------------------------

func init() {
	// Load configuration during initialization
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	openWeatherAPIURL = config.OpenWeatherAPIURL
	apiKey = config.APIKey

	// Log as JSON instead of the default ASCII formatter
	handlerLog.SetFormatter(&logrus.JSONFormatter{})
	handlerLog.SetOutput(os.Stdout)
	handlerLog.SetLevel(logrus.InfoLevel)
}

//----------------------------------------------------------------------------------------------------------------------

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// HTTP client for making requests
var client HTTPClient = &http.Client{}

//----------------------------------------------------------------------------------------------------------------------

// Config struct to hold configuration parameters
type Config struct {
	OpenWeatherAPIURL string `json:"openWeatherAPIURL"`
	APIKey            string `json:"apiKey"`
}

// LoadConfig reads the configuration from a JSON file
func loadConfig(filePath string) (*Config, error) {
	// Check if running in a test environment
	if os.Getenv("TEST_ENV") == "true" {
		return &Config{
			OpenWeatherAPIURL: "mockedAPIURL",
			APIKey:            "mockedAPIKey",
		}, nil
	}

	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	config := &Config{}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

//----------------------------------------------------------------------------------------------------------------------

// WeatherResponse represents the structure of the OpenWeatherMap API response
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

//----------------------------------------------------------------------------------------------------------------------

// WeatherHandler handles requests to the weather endpoint
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a weather request")
	vars := mux.Vars(r)
	lat, lon := GetCoordinates(vars)

	if lat == 0 || lon == 0 {
		handlerLog.Errorf("Invalid coordinates: lat : %v, lon %v", lat, lon)
		http.Error(w, "Invalid coordinates", http.StatusBadRequest)
		return
	}

	weatherData, err := getWeatherData(lat, lon)
	if err != nil {
		handlerLog.Errorf("Unable to fetch weather data: %v", err)
		http.Error(w, "Unable to fetch weather data", http.StatusInternalServerError)
		return
	}

	temperature := weatherData.Main.Temp
	weatherCondition := weatherData.Weather[0].Main
	description := weatherData.Weather[0].Description

	log.Printf("Weather information for coordinates (%f, %f): Temperature=%.2f°C, Condition=%s, Description=%s", lat, lon, temperature, weatherCondition, description)

	// Categorize weather conditions based on temperature
	weatherType := categorizeWeather(temperature)

	response := map[string]interface{}{
		"temperature":      temperature,
		"weatherCondition": weatherCondition,
		"description":      description,
		"weatherType":      weatherType,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		handlerLog.Errorf("Unable to encode response: %v", err)
		http.Error(w, "Unable to encode", http.StatusInternalServerError)
		return
	}
	log.Println("Weather request processed successfully")
}

//----------------------------------------------------------------------------------------------------------------------

// getCoordinates gets the latitude and longitude value
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

//----------------------------------------------------------------------------------------------------------------------

// getWeatherData fetches weather data from the OpenWeatherMap API
func getWeatherData(lat, lon float64) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", openWeatherAPIURL, lat, lon, apiKey)

	log.Printf("Fetching weather data from: %s", url)
	response, err := client.Get(url)
	if err != nil {
		handlerLog.Errorf("Error making HTTP request: %v", err)
		return nil, err
	}
	if response.StatusCode != 200 {
		handlerLog.Errorf("Received non-200 status code: %s", response.Status)
		return nil, fmt.Errorf(response.Status)
	}

	defer response.Body.Close()

	var weatherData WeatherResponse
	if err := json.NewDecoder(response.Body).Decode(&weatherData); err != nil {
		handlerLog.Errorf("Error decoding JSON response: %v", err)
		return nil, err
	}
	log.Println("Weather data fetched successfully")
	return &weatherData, nil
}

//----------------------------------------------------------------------------------------------------------------------

// categorizeWeather categorizes weather conditions based on temperature
func categorizeWeather(temperature float64) string {
	// Check if the temperature is less than or equal to 0
	if temperature <= 0 {
		return "Freezing"
	} else if temperature > 0 && temperature < 10 {
		// Check if the temperature is greater than 0 and less than 10
		return "Cold"
	} else if temperature >= 10 && temperature < 25 {
		// Check if the temperature is greater than or equal to 10 and less than 25
		return "Moderate"
	} else {
		// If the temperature is greater than or equal to 25
		return "Hot"
	}
}

//----------------------------------------------------------------------------------------------------------------------
