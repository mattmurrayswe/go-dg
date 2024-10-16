package main

import (
	"dg/controllers"
	"dg/db"
	"dg/middlewares"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func init() {
	db.InitializeDB()
}

func main() {
	http.HandleFunc("/", controllers.Status)

	// Auth routes
	http.HandleFunc("/login", controllers.Login)     //ṔOST
	http.HandleFunc("/refresh", controllers.Refresh) //GET
	http.HandleFunc("/user", controllers.CreateUser) //POST

	// Banner routes
	http.HandleFunc("/gen-banner", middlewares.Authenticate(controllers.GenBanner))         //POST
	http.HandleFunc("/logos", middlewares.Authenticate(controllers.GetLogos))               //GET
	http.HandleFunc("/tech-options", middlewares.Authenticate(controllers.ListTechOptions)) //GET

	// Car brands and models routes
	http.HandleFunc("/brands/", middlewares.Authenticate(controllers.ListBrands))     //GET
	http.HandleFunc("/brands/dg", middlewares.Authenticate(controllers.ListDgBrands)) //GET
	http.HandleFunc("/models", middlewares.Authenticate(controllers.ListModels))      //GET
	http.HandleFunc("/models/count", middlewares.Authenticate(controllers.ListModelsCount))      //GET
	http.HandleFunc("/models/detailed", middlewares.Authenticate(controllers.ListModelsDetailed))      //GET
	http.HandleFunc("/models/mixed", middlewares.Authenticate(controllers.ListModelsDetailed))      //GET
	http.HandleFunc("/versions", middlewares.Authenticate(controllers.ListVersions))      //GET
	http.HandleFunc("/versions/only", middlewares.Authenticate(controllers.ListVersionsOnly))      //GET
	http.HandleFunc("/versions/count", middlewares.Authenticate(controllers.ListModelsWithVersions))      //GET
	http.HandleFunc("/versions/detailed", middlewares.Authenticate(controllers.ListModelsWithVersions))      //GET
	http.HandleFunc("/versions/mixed", middlewares.Authenticate(controllers.ListModelsWithVersions))      //GET
	http.HandleFunc("/project", middlewares.Authenticate(controllers.CreateProject))     //POST
	http.HandleFunc("/wish-list", middlewares.Authenticate(controllers.WishList))     //POST
	// http.HandleFunc("/project-card/buy", middlewares.Authenticate(controllers.BuyProjectCard)) // POST

	fmt.Println("Server running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
