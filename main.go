package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//bannerdev
	http.HandleFunc("/gen-banner", genBanner)
	http.HandleFunc("/logos", getLogos) // Handles requests to get logos with URLs
	http.HandleFunc("/tech-options", listTechOptions) // Handles requests to get only tech names

	//digitalgarage
	http.HandleFunc("/brands", handleGetTechOptions)
	http.HandleFunc("/models", handleGetTechOptions)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
