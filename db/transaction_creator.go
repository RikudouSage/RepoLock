package db

import "database/sql"

type TransactionCreator interface {
	Create() (*sql.Tx, error)
}

func NewTransactionCreator(
	db *sql.DB,
) TransactionCreator {
	return &transactionCreator{
		db: db,
	}
}

type transactionCreator struct {
	db *sql.DB
}

func (receiver *transactionCreator) Create() (*sql.Tx, error) {
	return receiver.db.Begin()
}
