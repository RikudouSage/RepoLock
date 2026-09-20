package db

import (
	"context"
	"database/sql"
)

type QueryIssuer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
