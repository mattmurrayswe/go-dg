package controllers

import (
	"context"
	"dg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/plutov/paypal/v4"
)

// PayPalClient struct to manage PayPal client
type PayPalClient struct {
	Client *paypal.Client
}

// NewPayPalClient initializes the PayPal client
func NewPayPalClient() (*PayPalClient, error) {
	clientID := os.Getenv("PAYPAL_CLIENT_ID")
	secret := os.Getenv("PAYPAL_SECRET")
	if clientID == "" || secret == "" {
		return nil, fmt.Errorf("PayPal credentials are missing")
	}

	// Determine PayPal environment (sandbox or live)
	apiBase := paypal.APIBaseLive
	if os.Getenv("PAYPAL_ENV") == "sandbox" {
		apiBase = "https://api-m.sandbox.paypal.com" // Manually specify the sandbox URL
	}

	// Create PayPal client
	c, err := paypal.NewClient(clientID, secret, apiBase)
	if err != nil {
		return nil, err
	}
	c.SetLog(os.Stdout) // Set log for PayPal requests
	return &PayPalClient{Client: c}, nil
}

// CreatePayPalPayment creates an order in BRL
func CreatePayPalPayment(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Initialize PayPal client
	paypalClient, err := NewPayPalClient()
	if err != nil {
		http.Error(w, "Unable to create PayPal client", http.StatusInternalServerError)
		log.Printf("Error creating PayPal client: %v", err)
		return
	}

	// Define the order object with BRL currency
	purchaseUnits := []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Currency: "BRL",   // Currency set to Brazilian Real
				Value:    "50.00", // Set this dynamically as needed
			},
			Description: "Purchase of card from another user",
		},
	}

	// Application context for redirect URLs
	applicationContext := &paypal.ApplicationContext{
		ReturnURL: "http://localhost:8080/payment/success",
		CancelURL: "http://localhost:8080/payment/cancel",
	}

	// Create the order
	createdOrder, err := paypalClient.Client.CreateOrder(ctx, "CAPTURE", purchaseUnits, nil, applicationContext)
	if err != nil {
		http.Error(w, "Unable to create PayPal order", http.StatusInternalServerError)
		log.Printf("Error creating PayPal order: %v", err)
		return
	}

	// Redirect the user to the PayPal approval URL
	for _, link := range createdOrder.Links {
		if link.Rel == "approve" {
			http.Redirect(w, r, link.Href, http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Unable to find approval URL", http.StatusInternalServerError)
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
