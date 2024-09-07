package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Response structure for each service
type ServiceResponse struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

// Full response structure
type Response struct {
	Services []ServiceResponse `json:"services"`
}

func main() {
	http.HandleFunc("/", handleRequest)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
