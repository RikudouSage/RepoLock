package db

import (
	"context"
	"database/sql"
)

const transactionContextKey = "dbTransactionContext"

func WithTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, transactionContextKey, tx)
}

func NewTransactionContext(tx *sql.Tx) context.Context {
	return WithTransaction(context.Background(), tx)
}

func GetTransactionFromContext(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(transactionContextKey).(*sql.Tx)
	if !ok {
		return nil
	}

	return tx
}
