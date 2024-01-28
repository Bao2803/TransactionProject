package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to a method like
	// Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

type Transaction struct {
	gorm.Model
	Value    float64
	Note     string
	Sender   string
	Receiver string
}

type TransactionService struct {
	db *gorm.DB
}

// NewTransactionService Create a new TransactionService with a specified connectionInfo.
func NewTransactionService(connectionInfo string) (*TransactionService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)
	return &TransactionService{
		db: db,
	}, nil
}

// Close Closes the TransactionService database connection.
func (ts *TransactionService) Close() error {
	return ts.db.Close()
}

// AutoMigrate will attempt to automatically migrate the transactions table
func (ts *TransactionService) AutoMigrate() error {
	if err := ts.db.AutoMigrate(&Transaction{}).Error; err != nil {
		return err
	}
	return nil
}

// first will query using the provided gorm.DB, and it will get the first item
// returned and place it into dst. If nothing is found in the query, it will
// return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

// Create will create the provided user and back-fill data like
// the ID, CreatedAt, and UpdatedAt fields.
func (ts *TransactionService) Create(transaction *Transaction) error {
	return ts.db.Create(transaction).Error
}

// Update will update the provided user with all the data in the provided
// user object.
func (ts *TransactionService) Update(user *Transaction) error {
	return ts.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (ts *TransactionService) Delete(id uint) error {
	if id == 0 { // Go default uint is 0, Gorm will delete all rows if id is not provided
		return ErrInvalidID
	}
	user := Transaction{Model: gorm.Model{ID: id}}
	return ts.db.Delete(user).Error
}
