package controllers

import "transaction_project/models"

type Transactions struct {
	ts *models.TransactionService
}

func NewTransactions(ts *models.TransactionService) *Transactions {
	return &Transactions{
		ts: ts,
	}
}
