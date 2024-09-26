// banner.go
package controllers

import (
	"dg/db"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"

	"github.com/fogleman/gg"
)

type Technologies struct {
	Tech1 string `json:"tech1"`
	Tech2 string `json:"tech2"`
	Tech3 string `json:"tech3"`
	Tech4 string `json:"tech4"`
	Tech5 string `json:"tech5"`
}

func GenBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var techs Technologies
	err := json.NewDecoder(r.Body).Decode(&techs)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	img := generateBannerImage(techs)

	w.Header().Set("Content-Type", "image/png")

	err = png.Encode(w, img)
	if err != nil {
		http.Error(w, "Error generating image", http.StatusInternalServerError)
		return
	}
}

func generateBannerImage(techs Technologies) image.Image {
	const width = 1584
	const height = 396

	dc := gg.NewContext(width, height)
	dc.SetColor(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	dc.Clear()

	logoPaths := map[string]string{
		"Go":      "images/go.png",
		"PHP":     "images/php.png",
		"Laravel": "images/laravel.png",
		"React":   "images/react.png",
	}

	technologies := []string{techs.Tech1, techs.Tech2, techs.Tech3, techs.Tech4, techs.Tech5}
	maxHeight := 100.0
	margin := 50.0

	x := float64(width) - 100.0

	for _, tech := range technologies {
		logoPath, exists := logoPaths[tech]
		if exists {
			im, err := gg.LoadImage(logoPath)
			if err != nil {
				fmt.Println("Error loading logo:", err)
				continue
			}

			w := im.Bounds().Dx()
			h := im.Bounds().Dy()

			scale := maxHeight / float64(h)

			dc.Scale(scale, scale)
			dc.DrawImage(im, int(x/scale)-w, int((height-maxHeight)/scale))
			dc.Identity()

			x -= float64(w)*scale + margin
		}
	}

	return dc.Image()
}

type Tech struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

var availableTechs = []Tech{
	{"Go", "/images/go.png"},
	{"PHP", "/images/php.png"},
	{"Laravel", "/images/laravel.png"},
	{"React", "/images/react.png"},
	{"AWS", "/images/aws.png"},
}

func GetLogos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	host := r.Host
	techsWithFullURLs := make([]Tech, len(availableTechs))

	for i, tech := range availableTechs {
		techsWithFullURLs[i] = Tech{
			Name:     tech.Name,
			ImageURL: fmt.Sprintf("http://%s%s", host, tech.ImageURL),
		}
	}

	json.NewEncoder(w).Encode(techsWithFullURLs)
}

func ListTechOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.DB.Query("SELECT tech_name FROM tech_options")
	if err != nil {
		http.Error(w, "Unable to fetch tech options", http.StatusInternalServerError)
		log.Println("Error fetching tech options:", err)
		return
	}
	defer rows.Close()

	var techNames []string

	for rows.Next() {
		var techName string
		if err := rows.Scan(&techName); err != nil {
			http.Error(w, "Error scanning tech options", http.StatusInternalServerError)
			log.Println("Error scanning tech options:", err)
			return
		}
		techNames = append(techNames, techName)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating over tech options", http.StatusInternalServerError)
		log.Println("Error iterating over tech options:", err)
		return
	}

	if err := json.NewEncoder(w).Encode(techNames); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
	}
}
