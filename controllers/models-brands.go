package controllers

import (
	"encoding/json"
	"net/http"
)

var (
	brands = []string{
		"Acura", "Adamo", "Agrale", "Alfa Romeo", "Americar", "Asia", "Aston Martin", "Audi",
		"Austin-Healey", "Avallone", "Beach", "Bentley", "Bianco", "BMW", "BRM", "Bugre", "Bugway",
		"Buick", "BYD", "Cadillac", "Caoa Chery", "CBT", "Chamonix", "Cheda", "Chevrolet",
		"Chrysler", "Citroën", "Daewoo", "Daihatsu", "De Soto", "DKW-Vemag", "Dodge", "Edsel",
		"Effa", "Emis", "Engesa", "Envemo", "Farus", "Fercar Buggy", "Ferrari", "Fiat", "Ford",
		"Fyber", "Geely", "GMC", "Gurgel", "GWM", "Hafei", "Honda", "Hudson", "Hummer", "Hyundai",
		"Infiniti", "IVECO", "JAC", "Jaguar", "Jeep", "Jinbei", "Kia", "Lada", "Lamborghini",
		"Land Rover", "Lexus", "Lifan", "Lincoln", "Lotus", "Mahindra", "Marcopolo", "Maserati",
		"Mazda", "Mclaren", "Menon", "Mercedes-Benz", "Mercury", "MG", "Mini", "Mitsubishi", "Miura",
		"Mobby", "Morris", "MP Lafer", "Neta", "Nissan", "Opel", "PAG", "Peugeot", "Plymouth",
		"Pontiac", "Porsche", "Puma", "RAM", "Renault", "Rivian", "Rolls-Royce", "Saturn", "Seat",
		"Seres", "Shrlby", "Smart", "Ssangyong", "Studebaker", "Subaru", "Sunbeam", "Suzuki", "TAC",
		"Tesla", "Toyota", "Triumph", "Troller", "Volkswagen", "Volvo", "Wake", "Willys", "Willys Overland",
	}

	dgBrands = []string{
		"Alfa Romeo", "Aston Martin", "Audi", "BMW", "Bugatti", "Dodge", "Ferrari", "Jaguar", 
		"Lamborghini", "Land Rover", "Lexus", "Maserati", "Mclaren", "Mercedes-Benz", "Mini",
		"Porsche", "Rolls-Royce", "Volkswagen", 
	}

	brandsWithModels = map[string][]string{
		"Alfa Romeo": {"Giulia", "Stelvio", "Giulietta", "Stelvio Quadrifoglio", ""},
		"Aston Martin": {},
		"Audi": {},
		"BMW": {},
		"Bugatti": {},
	}
)

func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Application running on Kubernetes"}
	json.NewEncoder(w).Encode(response)
}

func ListBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brands)
}

func ListDgBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dgBrands)
}

func ListModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	brand := r.URL.Query().Get("brand")
	if brand != "" {
		if models, ok := brandsWithModels[brand]; ok {
			json.NewEncoder(w).Encode(models)
		} else {
			http.Error(w, "Brand not found", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Brand query parameter is missing", http.StatusBadRequest)
	}
}
