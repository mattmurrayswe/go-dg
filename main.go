package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//bannerdev
	http.HandleFunc("/gen-banner", handlePostRequest)
	http.HandleFunc("/logos", handleGetLogos) // Handles requests to get logos with URLs
	http.HandleFunc("/tech-options", handleGetTechOptions) // Handles requests to get only tech names

	//digitalgarage
	http.HandleFunc("/brands", handleGetTechOptions)
	http.HandleFunc("/models", handleGetTechOptions)
	
	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
