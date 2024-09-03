package main

import (
	"encoding/json"
	"net/http"
)

var (
	dgBrands = []string{
		"Alfa Romeo", "Audi", "Bentley", "BMW", "Bugatti", "Chevrolet", "Citroën", "Dodge", "Ferrari", "Ford", "Honda", "Jaguar", 
		"Kia", "Lamborghini", "Land Rover", "Lexus", "Maserati", "Mclaren", "Mercedes-Benz", "Mini", "Mitsubishi", "Nissan", 
		"Porsche", "RAM", "Subaru", "Toyota", "Volkswagen", "Volvo", 
	}

	allBrands = []string{
		"Alfa Romeo", "Americar", "Audi", "Bentley", "BMW", "Bugatti", "Bugway", "BYD", "Caoa Chery", "Chamonix", "Chevrolet", 
		"Chrysler", "Citroën", "Daihatsu", "Dodge", "Effa", "Engesa", "Ferrari", "Fiat", "Ford", "GMC", "GWM", "Honda", 
		"Hyundai", "Infiniti", "IVECO", "JAC", "Jaguar", "Jeep", "Kia", "Lamborghini", "Land Rover", "Lexus", "Lifan", 
		"Maserati", "Mclaren", "Mercedes-Benz", "Mini", "Mitsubishi", "Miura", "Mp Lafer", "Nissan", "Peugeot", "Plymouth", 
		"Porsche", "RAM", "Renault", "Rivian", "Seat", "Smart", "Ssangyong", "Studebaker", "Subaru", "Suzuki", "Tesla", 
		"Toyota", "Troller", "Volkswagen", "Volvo", "Willys", 
	}

	dgOmittedBrands = []string{
		"Americar", "Bugway", "BYD", "Caoa Chery", "Chamonix", "Chrysler", "Daihatsu", "Effa", "Engesa", "Fiat", "GMC", 
		"GWM", "Hyundai", "Infiniti", "IVECO", "JAC", "Jeep", "Lifan", "Miura", "MP Lafer", "Peugeot", "Plymouth", "Renault", 
		"Rivian", "Seat", "Smart", "Ssangyong", "Studebaker", "Suzuki", "Tesla", "Troller", "Willys", 
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
