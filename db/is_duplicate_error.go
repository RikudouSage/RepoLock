package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mattn/go-sqlite3"
)

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	if pgError, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgError.Code == "23505"
	}

	if sqliteError, ok := errors.AsType[sqlite3.Error](err); ok {
		return sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique ||
			sqliteError.ExtendedCode == sqlite3.ErrConstraintPrimaryKey
	}

	return false
}
