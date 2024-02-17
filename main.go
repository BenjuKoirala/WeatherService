package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/weather/{lat}/{lon}", weatherHandler).Methods("GET")

	port := 3000
	fmt.Printf("Server is running on :%d\n", port)
	http.Handle("/", router)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
