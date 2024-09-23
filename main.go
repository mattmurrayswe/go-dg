package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"  // Import PostgreSQL driver
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	http.HandleFunc("/", status) //GET

	http.HandleFunc("/login", login)
	http.HandleFunc("/user", createUser)
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

func createUser(w http.ResponseWriter, r *http.Request) {
    var username, password string
    username = r.FormValue("username")
    password = r.FormValue("password")

    // Validate input
    if username == "" || password == "" {
        http.Error(w, "Missing username or password", http.StatusBadRequest)
        return
    }

    // Hash the password before storing it
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    // Connection string for PostgreSQL
    connStr := "postgres://admin_dg:AmazingDGPass@123@192.168.49.2:30000/dg_go?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        log.Println("Error connecting to the database:", err)
        return
    }
    defer db.Close()

    // Insert the new user into the database (use $1, $2 for parameterized queries in PostgreSQL)
    _, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashedPassword))
    if err != nil {
        http.Error(w, "Could not create user", http.StatusInternalServerError)
        log.Println("Database error:", err) // Log error for further investigation
        return
    }

    // Respond with a success message
    json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func login(w http.ResponseWriter, r *http.Request) {
    var username, password string
    username = r.FormValue("username")
    password = r.FormValue("password")

    // Validate input
    if username == "" || password == "" {
        http.Error(w, "Missing username or password", http.StatusBadRequest)
        return
    }

    // Connect to the PostgreSQL database
    connStr := "postgres://admin_dg:AmazingDGPass@123@192.168.49.2:30000/dg_go?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        log.Println("Error connecting to the database:", err)
        return
    }
    defer db.Close()

    // Query the database for the user's hashed password
    var hashedPassword string
    err = db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
    if err == sql.ErrNoRows {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "Database query error", http.StatusInternalServerError)
        log.Println("Error querying the database:", err)
        return
    }

    // Compare the hashed password with the provided password
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        // Password does not match
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Create the JWT claims, which include the username and expiry time
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
