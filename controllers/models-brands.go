package controllers

import (
	"database/sql"
	"dg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Brand struct {
	BrandName string `json:"brand_name"`
}

type Model struct {
	ID        int    `json:"id"`
	BrandName string `json:"brand_name"`
	ModelName string `json:"model_name"`
}

type Version struct {
	ID          int    `json:"id"`
	BrandName   string `json:"brand_name"`
	ModelName   string `json:"model_name"`
	VersionName string `json:"version_name"`
}

type BrandModels struct {
	Brand  string   `json:"brand"`
	Logo   string   `json:"logo"`
	Site   string   `json:"data_source_website"`
	Models []string `json:"models"`
}

type Brands struct {
	Brands []string `json:"brands"`
}

func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Application running on Kubernetes"}
	json.NewEncoder(w).Encode(response)
}

func ListBrands(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.DB.Query(`
		SELECT brand_name
		FROM brands
	`)

	var brands Brands
	for rows.Next() {
		var brand Brand
		rows.Scan(&brand.BrandName)
		brands.Brands = append(brands.Brands, brand.BrandName)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brands)
}

func ListDgBrands(w http.ResponseWriter, r *http.Request) {

	rows, _ := db.DB.Query(`
		SELECT brand_name
		FROM dg_brands
	`)

	var brands []string
	for rows.Next() {
		var brand string
		rows.Scan(&brand)
		brands = append(brands, brand)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brands)
}

type BrandModelsOnly struct {
	BrandName string `json:"brand"`
	ModelNames []string `json:models`
}

func ListModels(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Listing Models")

	rows, _ := db.DB.Query(`
		SELECT *
		FROM models
	`)

	var models []BrandModelsOnly
	var brandIteration string
	for rows.Next() {

		var model, brand string
		var brandModels BrandModelsOnly
		rows.Scan(&model, &brand)
		if brandIteration != brand {
			brandModels.BrandName = brand
		}
		brandModels.ModelNames = append(brandModels.ModelNames, model)

		models = append(models, )
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

func ListModelsCount(w http.ResponseWriter, r *http.Request) {

	var count int
	db.DB.QueryRow(`
		SELECT COUNT (*)
		FROM models
	`).Scan(&count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)
	// json.NewEncoder(w).Encode(fmt.Sprintf("%d models found", count))
}

func ListVersions(w http.ResponseWriter, r *http.Request) {

	rows, _ := db.DB.Query(`
		SELECT *
		FROM versions
	`)

	var versions []Version
	for rows.Next() {
		var version Version
		rows.Scan(&version.ID, &version.BrandName, &version.ModelName, &version.VersionName)
		versions = append(versions, version)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}

func ListModelsDetailed(w http.ResponseWriter, r *http.Request) {
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
				Brand:  brandName,
				Site:   siteValue,
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
		Brand              string   `json:"brand"`
		Site               string   `json:"site"`
		Logo               string   `json:"logo"`
		Models             []string `json:"models"`            // List of concatenated model and version
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
				Brand:              brandName,
				Site:               site,
				Logo:               logoValue,
				Models:             []string{}, // Initialize an empty slice for models with versions
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

type VersionsOnly struct {
	VersionsConcat []string `json:"versions_concat"`
}

func ListVersionsOnly(w http.ResponseWriter, r *http.Request) {

	rows, _ := db.DB.Query(`
		SELECT * FROM
		versions
	`)

	var versionOnly VersionsOnly
	for rows.Next() {
		var version Version
		rows.Scan(&version.ID, &version.BrandName, &version.ModelName, &version.VersionName)
		fmt.Print(version.BrandName + version.ModelName + version.VersionName)
		versionOnly.VersionsConcat = append(versionOnly.VersionsConcat, version.BrandName+` `+version.ModelName+` `+version.VersionName)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionOnly)
}
