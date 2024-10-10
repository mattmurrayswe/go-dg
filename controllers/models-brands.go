package controllers

import (
	"database/sql"
	"dg/db"
	"encoding/json"
	"log"
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
		"Alfa Romeo":   {"Giulia", "Stelvio", "Giulietta", "Stelvio Quadrifoglio", ""},
		"Aston Martin": {},
		"Audi":         {},
		"BMW":          {},
		"Bugatti":      {},
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

	// Query models with their brand and site
	rows, err := db.DB.Query(`
		SELECT models.model_name, models.brand_name, dg_brands.site, dg_brands.logo_url 
		FROM models
		INNER JOIN dg_brands ON models.brand_name = dg_brands.brand_name
	`)
	if err != nil {
		http.Error(w, "Unable to fetch models", http.StatusInternalServerError)
		log.Printf("Error fetching models: %v", err)
		return
	}
	defer rows.Close()

	type BrandModels struct {
		Brand              string              `json:"brand"`
		Logo               string              `json:"logo"`
		Site               string              `json:"data_source_website"`
		Models             []string            `json:"models"`
	}

	// Use a map to group models by brand and store the source site
	brandModelsMap := make(map[string]BrandModels)

	// Process rows of models
	for rows.Next() {
		var modelName, brandName string
		var logo, site sql.NullString

		if err := rows.Scan(&modelName, &brandName, &site, &logo); err != nil {
			http.Error(w, "Error scanning models", http.StatusInternalServerError)
			log.Printf("Error scanning models: %v", err)
			return
		}

		// Convert sql.NullString to string, use empty string if null
		logoValue := ""
		if logo.Valid {
			logoValue = logo.String
		}

		// Convert sql.NullString to string, use empty string if null
		siteValue := ""
		if site.Valid {
			siteValue = site.String
		}

		// Check if the brand already exists in the map
		if brandModel, exists := brandModelsMap[brandName]; exists {
			brandModel.Models = append(brandModel.Models, modelName)
			brandModelsMap[brandName] = brandModel
		} else {
			// Create a new BrandModels entry for this brand
			brandModelsMap[brandName] = BrandModels{
				Brand:              brandName,
				Site:               siteValue,
				Logo:               logoValue,
				Models:             []string{modelName},
			}
		}
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error fetching data from database", http.StatusInternalServerError)
		log.Printf("Error after row iteration: %v", err)
		return
	}

	// Convert map to slice of BrandModels
	var result []BrandModels
	for _, brandModel := range brandModelsMap {
		result = append(result, brandModel)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func ListModelsWithVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query models with their brand and site
	rows, err := db.DB.Query(`
		SELECT models.model_name, models.brand_name, dg_brands.site, dg_brands.logo_url 
		FROM models
		INNER JOIN dg_brands ON models.brand_name = dg_brands.brand_name
	`)
	if err != nil {
		http.Error(w, "Unable to fetch models", http.StatusInternalServerError)
		log.Printf("Error fetching models: %v", err)
		return
	}
	defer rows.Close()

	// Query versions associated with models
	rowsVersions, err := db.DB.Query(`
		SELECT versions.brand_name, versions.model_name, versions.version_name 
		FROM versions
	`)
	if err != nil {
		http.Error(w, "Unable to fetch versions", http.StatusInternalServerError)
		log.Printf("Error fetching versions: %v", err)
		return
	}
	defer rowsVersions.Close()

	type BrandModels struct {
		Brand  string   `json:"brand"`
		Site   string   `json:"site"`
		Logo   string   `json:"logo"`
		Models []string `json:"models"` // List of concatenated model and version
		ModelsWithVersions []string `json:"models_w_versions"` // List of concatenated model and version
	}

	// Use a map to group models by brand and store the source site
	brandModelsMap := make(map[string]BrandModels)

	// Process rows of models
	for rows.Next() {
		var modelName, brandName, site string
		var logo sql.NullString

		if err := rows.Scan(&modelName, &brandName, &site, &logo); err != nil {
			http.Error(w, "Error scanning models", http.StatusInternalServerError)
			log.Printf("Error scanning models: %v", err)
			return
		}

		// Convert sql.NullString to string, use empty string if null
		logoValue := ""
		if logo.Valid {
			logoValue = logo.String
		}

		// Check if the brand already exists in the map
		if brandModel, exists := brandModelsMap[brandName]; exists {
			brandModelsMap[brandName] = brandModel
		} else {
			// Create a new BrandModels entry for this brand
			brandModelsMap[brandName] = BrandModels{
				Brand:  brandName,
				Site:   site,
				Logo:   logoValue,
				Models: []string{}, // Initialize an empty slice for models with versions
				ModelsWithVersions: []string{}, // Initialize an empty slice for models with versions
			}
		}
	}

	// Process rows of versions
	for rowsVersions.Next() {
		var versionName, modelName, brandName string

		if err := rowsVersions.Scan(&brandName, &modelName, &versionName); err != nil {
			http.Error(w, "Error scanning versions", http.StatusInternalServerError)
			log.Printf("Error scanning versions: %v", err)
			return
		}

		// Create full model-version entry (e.g., "Stelvio Sprint")
		fullModelName := modelName
		if versionName != "" {
			fullModelName = modelName + " " + versionName
		}

		// Add the full model name to the list
		if brandModel, exists := brandModelsMap[brandName]; exists {
			brandModel.ModelsWithVersions = append(brandModel.ModelsWithVersions, fullModelName)
			brandModelsMap[brandName] = brandModel
		}
	}

	// Ensure all models have at least their base name if no version is present
	for brandName, brandModel := range brandModelsMap {
		if len(brandModel.ModelsWithVersions) == 0 {
			// If no versions exist for a model, just add the model name
			brandModel.ModelsWithVersions = append(brandModel.ModelsWithVersions, brandName)
		}
		brandModelsMap[brandName] = brandModel
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error fetching data from database", http.StatusInternalServerError)
		log.Printf("Error after row iteration: %v", err)
		return
	}

	// Convert map to slice of BrandModels
	var result []BrandModels
	for _, brandModel := range brandModelsMap {
		result = append(result, brandModel)
	}

	// Return the result as a JSON response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON response: %v", err)
	}
}
