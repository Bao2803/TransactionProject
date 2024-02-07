package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"transaction_project/controllers"
	"transaction_project/models"
)

// TODO: move to config file
const (
	host     = ""
	port     = 5432
	user     = ""
	password = ""
	dbname   = "postgres"
)

const defaultPort = "3000"

func main() {
	var err error
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = defaultPort
	}

	// Initiate database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	db.LogMode(true)
	log.Printf("Database connection established!")

	// Initiate services and AutoMigrate
	transService, err := models.NewTransactionService(db)
	if err != nil && transService.AutoMigrate() != nil {
		panic(err)
	}
	userService, err := models.NewUserService(db)
	if err != nil && userService.AutoMigrate() != nil {
		panic(err)
	}

	// Initiate controllers
	transController := controllers.NewTransactionController(transService)
	userController := controllers.NewUserController(userService)
	graphController := controllers.NewGraphQL(transController, userController)

	// Add handler and start server
	http.Handle("/graph", graphController.NewHandler())
	log.Printf("Connect to http://localhost:%s/graph for GraphQL playground", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
