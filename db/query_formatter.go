package db

import (
	"github.com/jmoiron/sqlx"
	"go.chrastecky.dev/repolock/config/data"
)

type QueryFormatter interface {
	FormatQuery(query string) string
}

func NewQueryFormatter(config *data.GlobalConfig) QueryFormatter {
	return &queryFormatter{dbType: config.DatabaseType}
}

type queryFormatter struct {
	dbType data.DatabaseType
}

func (receiver *queryFormatter) FormatQuery(query string) string {
	if receiver.dbType == data.DatabaseTypeSQLite {
		return sqlx.Rebind(sqlx.NAMED, query)
	}

	return sqlx.Rebind(sqlx.DOLLAR, query)
}
