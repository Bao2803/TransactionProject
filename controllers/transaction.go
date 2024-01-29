package controllers

import "transaction_project/models"

type Transaction struct {
	transService *models.TransactionService
}

// NewTransactionController create a new Transaction controller using the provided TransactionService
func NewTransactionController(transService *models.TransactionService) *Transaction {
	return &Transaction{
		transService: transService,
	}
}

// NewModel create a new models.Transaction and then add it to the database using the models.TransactionService
// Argument has type interface{} if it is not required
func (tC *Transaction) NewModel(value float64, note interface{}, sender, receiver string) (*models.Transaction, error) {
	var newNote string
	if note != nil {
		newNote = note.(string)
	}

	newTransaction := &models.Transaction{
		Value:    value,
		Note:     newNote,
		Sender:   sender,
		Receiver: receiver,
	}
	return newTransaction, tC.transService.Create(newTransaction)
}

// UpdateModel update an existed models.Transaction using the models.TransactionService
// Argument has type interface{} if it is not required
func (tC *Transaction) UpdateModel(id uint, value, note, sender, receiver interface{}) (*models.Transaction, error) {
	// Get the exist model
	transaction, err := tC.transService.ReadByID(id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	// Update exist model value
	updateValue := transaction.Value
	if value != nil {
		updateValue = value.(float64)
	}
	transaction.Value = updateValue

	// Update exist model note
	updateNote := transaction.Note
	if note != nil {
		updateNote = note.(string)
	}
	transaction.Note = updateNote

	// Update exist model sender
	updateSender := transaction.Sender
	if sender != nil {
		updateSender = sender.(string)
	}
	transaction.Sender = updateSender

	// Update exist model receiver
	updateReceiver := transaction.Receiver
	if receiver != nil {
		updateReceiver = receiver.(string)
	}
	transaction.Receiver = updateReceiver

	return transaction, tC.transService.Update(transaction)
}
