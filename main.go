package main

import (
	"WeatherService/weather"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

//----------------------------------------------------------------------------------------------------------------------

// Logrus logger instance
var log = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter
	log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above
	log.SetLevel(logrus.InfoLevel)
}

//----------------------------------------------------------------------------------------------------------------------

func main() {
	// Create a new router
	router := mux.NewRouter()
	setupEndpoints(router)

	port := 3000

	log.Infof("Server is running on :%d", port)

	// Handle requests using the router
	http.Handle("/", router)

	// Start the HTTP server and handle errors if any
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Errorf("Error starting server: %v", err)
	}
}

//----------------------------------------------------------------------------------------------------------------------

// setupEndpoints configures the routes for the application
func setupEndpoints(router *mux.Router) {
	// Define a route for the weather endpoint and specify the corresponding handler
	router.HandleFunc("/weather/{lat}/{lon}", weather.WeatherHandler).Methods("GET")
}

//----------------------------------------------------------------------------------------------------------------------
