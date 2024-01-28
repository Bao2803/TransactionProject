package main

import (
	"fmt"
	"net/http"
	"transaction_project/models"
)

// TODO: move to config file
const (
	host     = "transaction-project.postgres.database.azure.com"
	port     = 5432
	user     = "transaction"
	password = "222Wwood@"
)

func graph(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello, world!"))
	if err != nil {
		return
	}
}

func main() {
	// Create a DB connection string and then use it to create our model services.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
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

	// Initiate controllers
	//transactionC := controllers.NewTransactions(ts)

	http.HandleFunc("/graph", graph)
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
