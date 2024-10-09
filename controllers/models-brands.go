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
		Brand  string   `json:"brand"`
		Site   string   `json:"site"`
		Logo   string   `json:"logo"`
		Models []string `json:"models"`
	}

	// Use a map to group models by brand and store the source site
	brandModelsMap := make(map[string]BrandModels)

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
			brandModel.Models = append(brandModel.Models, modelName)
			brandModelsMap[brandName] = brandModel
		} else {
			// Create a new BrandModels entry for this brand
			brandModelsMap[brandName] = BrandModels{
				Brand:  brandName,
				Site:   site,
				Logo:   logoValue,
				Models: []string{modelName},
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


func ListModels1(w http.ResponseWriter, r *http.Request) {
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

func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type Project struct {
		OwnerID     *int    `json:"owner_id"`
		Brand       *string `json:"brand"`
		Model       *string `json:"model"`
		Year        *int    `json:"year"`
		CardPrice   *int    `json:"card_price"`
		ProjectName *string `json:"project_name"`
		Photo       *string `json:"photo"`
		HorsePowers *int    `json:"horse_powers"`
		DGP         *int    `json:"dgp"`
		Rarity      *string `json:"rarity"`
	}

	var project Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO projects (owner_id, brand, model, year, card_price, project_name, photo, horse_powers, dgp, rarity)
                     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = db.DB.Exec(sqlStatement, project.OwnerID, project.Brand, project.Model, project.Year, project.CardPrice,
		project.ProjectName, project.Photo, project.HorsePowers, project.DGP, project.Rarity)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project created successfully"})
}
