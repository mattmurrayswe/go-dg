package main

import (
	"dg/controllers"
	"dg/db"
	"dg/middlewares"
	"fmt"
	"log"
	"net/http"
	_ "github.com/lib/pq"
)

func init() {
    db.InitializeDB()
}

func main() {
	http.HandleFunc("/", controllers.Status)

	// Auth routes
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/refresh", controllers.Refresh)
	http.HandleFunc("/user", controllers.CreateUser)

	// Banner routes
	http.HandleFunc("/gen-banner", middlewares.Authenticate(controllers.GenBanner))
	http.HandleFunc("/logos", middlewares.Authenticate(controllers.GetLogos))
	http.HandleFunc("/tech-options", middlewares.Authenticate(controllers.ListTechOptions))

	// Car brands and models routes
	http.HandleFunc("/brands/", middlewares.Authenticate(controllers.ListBrands))
	http.HandleFunc("/brands/dg", middlewares.Authenticate(controllers.ListDgBrands))
	http.HandleFunc("/models", middlewares.Authenticate(controllers.ListModels))

	fmt.Println("Server running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
