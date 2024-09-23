package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq" // PostgreSQL driver
    "github.com/joho/godotenv"
)

var DB *sql.DB

func init() {
    // Load the .env file
    fmt.Println("2")
    err := godotenv.Load("./.env")
    if err != nil {
        fmt.Println("3")
        log.Fatalf("Error loading .env file: %v", err)
    }
}

// InitializeDB loads the .env file and establishes the DB connection
func InitializeDB() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    fmt.Println("5")
    // Get connection details from environment variables
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    // Construct the PostgreSQL connection string
    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

    // Open the database connection
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database connection: %v", err)
    }

    // Verify the connection
    err = DB.Ping()
    if err != nil {
        log.Fatalf("Error verifying database connection: %v", err)
    }

    log.Println("Database connection established")
}
