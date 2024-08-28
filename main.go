package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//bannerdev
	http.HandleFunc("/gen-banner", genBanner)         //POST
	http.HandleFunc("/logos", getLogos)               //GET
	http.HandleFunc("/tech-options", listTechOptions) //GET

	//cars-models-brands
	http.HandleFunc("/brands/all", listAllBrands)           //GET
	http.HandleFunc("/brands", listDgBrands)                //GET
	http.HandleFunc("/brands/omitted", listDgOmmitedBrands) //GET
	http.HandleFunc("/models", listModels)                  //GET

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
