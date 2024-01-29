package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"transaction_project/controllers"
	"transaction_project/models"
)

// TODO: move to config file
const (
	host     = "transaction-project.postgres.database.azure.com"
	port     = 5432
	user     = "transaction"
	password = "222Wwood"
	dbname   = "postgres"
)

const defaultPort = "3000"

func main() {
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = defaultPort
	}

	// Create a DB connection string and then use it to create our model services.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	ts, err := models.NewTransactionService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func(ts *models.TransactionService) {
		err := ts.Close()
		if err != nil {
			panic(err)
		}
	}(ts)
	err = ts.AutoMigrate()
	if err != nil {
		panic(err)
	}
	log.Printf("Database connection established!")

	// Initiate controllers
	transactionC := controllers.NewTransactionController(ts)
	graphqlC := controllers.NewGraphQL(transactionC) // Binding controllers to the graphql controller

	// Add handler and start server
	http.Handle("/graph", graphqlC.NewHandler())
	log.Printf("Connect to http://localhost:%s/graph for GraphQL playground", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
