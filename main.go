package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", status) //GET

	//bannerdev
	http.HandleFunc("/gen-banner", genBanner)         //POST
	http.HandleFunc("/logos", getLogos)               //GET
	http.HandleFunc("/tech-options", listTechOptions) //GET

	//cars-models-brands
	http.HandleFunc("/brands/", listBrands)                    //GET
	http.HandleFunc("/brands/dg", listDgBrands)                //GET
	http.HandleFunc("/models", listModels)                     //GET

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
