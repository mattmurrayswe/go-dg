package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	http.HandleFunc("/", status) //GET

	http.HandleFunc("/login", login)
	http.HandleFunc("/refresh", refresh) // New route for refreshing tokens
	http.HandleFunc("/protected", authenticate(protectedEndpoint))

	//bannerdev
	http.HandleFunc("/gen-banner", authenticate(genBanner))         //POST
	http.HandleFunc("/logos", authenticate(getLogos))               //GET
	http.HandleFunc("/tech-options", authenticate(listTechOptions)) //GET

	//cars-models-brands
	http.HandleFunc("/brands/", authenticate(listBrands))     //GET
	http.HandleFunc("/brands/dg", authenticate(listDgBrands)) //GET
	http.HandleFunc("/models", authenticate(listModels))      //GET

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the protected endpoint!"))
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		// Parse the token and validate
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid token", http.StatusBadRequest)
			}
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid, continue
		next(w, r)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	var username, password string
	username = r.FormValue("username")
	password = r.FormValue("password")

	// Validate credentials (simplified)
	if username != "admin" || password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create the JWT claims, which includes the username and expiry time
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func refresh(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if token is close to expiry
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		http.Error(w, "Token is still valid", http.StatusBadRequest)
		return
	}

	// Issue a new token
	expirationTime := time.Now().Add(15 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = newToken.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not refresh token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
