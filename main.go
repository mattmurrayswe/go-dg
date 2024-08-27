package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Tech struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

var availableTechs = []Tech{
	{"Go", "/images/go.png"},
	{"PHP", "/images/php.png"},
	{"Laravel", "/images/laravel.png"},
	{"React", "/images/react.png"},
	{"AWS", "/images/aws.png"},
}

func main() {
	http.HandleFunc("/gen-banner", handlePostRequest)
	http.HandleFunc("/logos", handleGetLogos) // Handles requests to get logos with URLs
	http.HandleFunc("/tech-options", handleGetTechOptions) // Handles requests to get only tech names
	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetLogos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	host := r.Host
	techsWithFullURLs := make([]Tech, len(availableTechs))

	for i, tech := range availableTechs {
		techsWithFullURLs[i] = Tech{
			Name:     tech.Name,
			ImageURL: fmt.Sprintf("http://%s%s", host, tech.ImageURL),
		}
	}

	json.NewEncoder(w).Encode(techsWithFullURLs)
}

func handleGetTechOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	techNames := make([]string, len(availableTechs))
	for i, tech := range availableTechs {
		techNames[i] = tech.Name
	}

	json.NewEncoder(w).Encode(techNames)
}
