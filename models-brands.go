package main

import (
	"encoding/json"
	"net/http"
)

var (
	dgBrands = []string{
		"BMW", "Chevrolet", "Dodge", "Ford", "Honda", "Hyundai", "Jaguar", "Kia", "Land Rover",
		"Lexus", "Maserati", "Mercedez-Benz", "Mini", "Mitsubishi", "Nissan", "Peugeot", "Porsche",
		"RAM", "Renault", "Subaru", "Suzuki", "Toyota", "Volkswagen", "Volvo",
	}

	allBrands = []string{
		"BMW",  "BYD", "Caoa Chery", "Chevrolet", "Chrysler", "Citroën", "Dodge", "Fiat", "Ford",
		"GMC", "GWM", "Honda", "Hummer", "Hyundai", "Iveco", "JAC", "Jaguar", "Jeep", "Kia",
		"Land Rover", "Lexus", "Lifan", "Maserati", "Mercedez-Benz", "Mini", "Mitsubishi", "Nissan",
		"Peugeot", "Porsche", "RAM", "Renault", "Smart", "Subaru", "Suzuki", "Toyota", "Troller",
		"Volkswagen", "Volvo",
	}

	dgOmittedBrands = []string{
		"GMC", "GWM", "Hummer", "Iveco", "JAC", "Jeep", "Lifan", "Smart", "Troller",
	}

	models = []string{"Corolla", "Mustang", "X5", "C-Class"}
)

func listDgBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dgBrands)
}

func listDgOmmitedBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dgOmittedBrands)
}

func listAllBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allBrands)
}

func listModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}
