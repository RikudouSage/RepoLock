package repo

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/db"
)

type DefaultRepository[TEntity any] interface {
	FindByID(id uuid.UUID) (*TEntity, error)
	Find(options ...FindOption) ([]*TEntity, error)
}

type defaultRepository[TEntity any] struct {
	db             *sql.DB
	tableName      string
	queryFormatter db.QueryFormatter
	mapper         db.Mapper
}

func newDefaultRepository[TEntity any](
	db *sql.DB,
	tableName string,
	formatter db.QueryFormatter,
	mapper db.Mapper,
) *defaultRepository[TEntity] {
	return &defaultRepository[TEntity]{
		db:             db,
		tableName:      tableName,
		queryFormatter: formatter,
		mapper:         mapper,
	}
}

func (receiver *defaultRepository[TEntity]) FindByID(id uuid.UUID) (*TEntity, error) {
	result, err := receiver.Find(
		WithWhere("id = ?", id),
		WithLimit(1),
	)

	if err != nil {
		return nil, fmt.Errorf("failed fetching entity by id '%s': %w", id, err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}

func (receiver *defaultRepository[TEntity]) Find(options ...FindOption) ([]*TEntity, error) {
	result := make([]*TEntity, 0)
	query, bind := receiver.createQuery(options)
	query = receiver.queryFormatter.FormatQuery(query)
	rows, err := receiver.db.Query(query, bind...)
	if err != nil {
		var zero TEntity
		return nil, fmt.Errorf("failed fetching entities of type %T: %w", zero, err)
	}
	defer rows.Close()

	for rows.Next() {
		var item TEntity
		err = receiver.mapper.Map(rows, &item)
		if err != nil {
			return nil, fmt.Errorf("failed mapping entity of type %T: %w", item, err)
		}

		result = append(result, &item)
	}

	return result, nil
}

func (receiver *defaultRepository[TEntity]) createQuery(options []FindOption) (query string, bind []any) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("select * from ")
	queryBuilder.WriteString(receiver.tableName)

	optionsHolder := &findOptions{}
	for _, option := range options {
		option(optionsHolder)
	}

	if len(optionsHolder.where) != 0 {
		queryBuilder.WriteString(" where ")
		subQueries := make([]string, 0, len(optionsHolder.where))
		for _, where := range optionsHolder.where {
			subQueries = append(subQueries, "("+where+")")
		}
		queryBuilder.WriteString(strings.Join(subQueries, " and "))
	}

	if optionsHolder.orderByDirection != "" && optionsHolder.orderByField != "" {
		queryBuilder.WriteString(" order by " + optionsHolder.orderByField + " " + optionsHolder.orderByDirection)
	}

	if optionsHolder.limit != 0 {
		queryBuilder.WriteString(" limit " + strconv.FormatUint(uint64(optionsHolder.limit), 10))
	}

	return queryBuilder.String(), optionsHolder.bindValues
}
