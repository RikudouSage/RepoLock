package repo

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.chrastecky.dev/repolock/db"
)

type queryType string

const (
	queryTypeSelect queryType = "select"
	queryTypeDelete queryType = "delete"
)

type DefaultRepository[TEntity any] interface {
	FindByID(ctx context.Context, id uuid.UUID) (*TEntity, error)
	Find(ctx context.Context, options ...FindOption) ([]*TEntity, error)
	Create(ctx context.Context, entity *TEntity) error
	Delete(ctx context.Context, options ...FindOption) error
	Update(ctx context.Context, entity *TEntity) error
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

func (receiver *defaultRepository[TEntity]) FindByID(ctx context.Context, id uuid.UUID) (*TEntity, error) {
	result, err := receiver.Find(
		ctx,
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

func (receiver *defaultRepository[TEntity]) Find(ctx context.Context, options ...FindOption) ([]*TEntity, error) {
	database := receiver.getQueryIssuer(ctx)

	result := make([]*TEntity, 0)
	query, bind := receiver.createQuery(queryTypeSelect, options)
	query = receiver.queryFormatter.FormatQuery(query)
	rows, err := database.QueryContext(ctx, query, bind...)
	if err != nil {
		var zero TEntity
		return nil, fmt.Errorf("failed fetching entities of type %T: %w", zero, err)
	}
	defer rows.Close()

	for rows.Next() {
		var item TEntity
		err = receiver.mapper.MapOntoStruct(rows, &item)
		if err != nil {
			return nil, fmt.Errorf("failed mapping entity of type %T: %w", item, err)
		}

		result = append(result, &item)
	}

	if err := rows.Err(); err != nil {
		var zero TEntity
		return nil, fmt.Errorf("failed iterating entities of type %T: %w", zero, err)
	}

	return result, nil
}

func (receiver *defaultRepository[TEntity]) Create(ctx context.Context, entity *TEntity) error {
	database := receiver.getQueryExecutor(ctx)

	if err := receiver.createID(entity); err != nil {
		return fmt.Errorf("failed creating a new ID for entity: %w", err)
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("insert into ")
	queryBuilder.WriteString(receiver.tableName)
	queryBuilder.WriteString(" (")

	columns, values, err := receiver.mapper.MapFromStruct(entity)
	if err != nil {
		return fmt.Errorf("failed mapping entity of type %T: %w", entity, err)
	}

	queryBuilder.WriteString(strings.Join(columns, ", "))
	queryBuilder.WriteString(") values (")
	queryBuilder.WriteString(strings.Join(lo.RepeatBy(len(values), func(_ int) string {
		return "?"
	}), ", "))
	queryBuilder.WriteString(")")

	query := receiver.queryFormatter.FormatQuery(queryBuilder.String())
	if _, err = database.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("failed persisting entity %T: %w", entity, err)
	}

	return nil
}

func (receiver *defaultRepository[TEntity]) Delete(ctx context.Context, options ...FindOption) error {
	database := receiver.getQueryExecutor(ctx)

	query, bind := receiver.createQuery(queryTypeDelete, options)
	query = receiver.queryFormatter.FormatQuery(query)

	if _, err := database.ExecContext(ctx, query, bind...); err != nil {
		return fmt.Errorf("failed deleting entities: %w", err)
	}

	return nil
}

func (receiver *defaultRepository[TEntity]) Update(ctx context.Context, entity *TEntity) error {
	database := receiver.getQueryExecutor(ctx)

	var queryBuilder strings.Builder
	queryBuilder.WriteString("update ")
	queryBuilder.WriteString(receiver.tableName)
	queryBuilder.WriteString(" set ")

	columns, values, err := receiver.mapper.MapFromStruct(entity)
	if err != nil {
		return fmt.Errorf("failed mapping entity of type %T: %w", entity, err)
	}

	idColumnIndex := slices.Index(columns, "id")
	if idColumnIndex != -1 {
		columns = lo.DropByIndex(columns, idColumnIndex)
		values = lo.DropByIndex(values, idColumnIndex)
	}

	parts := make([]string, len(columns))
	for i, column := range columns {
		parts[i] = fmt.Sprintf("%s = ?", column)
	}
	queryBuilder.WriteString(strings.Join(parts, ", "))
	queryBuilder.WriteString(" where id = ?")

	values = append(values, lo.Must(receiver.getID(entity)))

	query := receiver.queryFormatter.FormatQuery(queryBuilder.String())
	if _, err = database.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("failed updating entity %T: %w", entity, err)
	}

	return nil
}

func (receiver *defaultRepository[TEntity]) getID(entity *TEntity) (uuid.UUID, error) {
	ref := reflect.ValueOf(entity).Elem()
	if ref.Kind() != reflect.Struct {
		return uuid.Nil, fmt.Errorf("entity is not a struct type: %T", entity)
	}

	field, ok := ref.Type().FieldByName("ID")
	if !ok {
		return uuid.Nil, fmt.Errorf("entity %T has no field 'ID'", entity)
	}

	if !field.Type.AssignableTo(reflect.TypeFor[uuid.UUID]()) {
		return uuid.Nil, fmt.Errorf("entity %T has an invalid type for its ID: %T", entity, ref.FieldByName(field.Name).Interface())
	}

	return ref.FieldByName(field.Name).Interface().(uuid.UUID), nil
}

func (receiver *defaultRepository[TEntity]) createID(entity *TEntity) error {
	ref := reflect.ValueOf(entity).Elem()
	if ref.Kind() != reflect.Struct {
		return fmt.Errorf("entity is not a struct type: %T", entity)
	}

	field, ok := ref.Type().FieldByName("ID")
	if !ok {
		return fmt.Errorf("entity %T has no field 'ID'", entity)
	}

	if !field.Type.AssignableTo(reflect.TypeFor[uuid.UUID]()) {
		return fmt.Errorf("entity %T has an invalid type for its ID: %T", entity, ref.FieldByName(field.Name).Interface())
	}

	ref.FieldByName(field.Name).Set(reflect.ValueOf(lo.Must(uuid.NewV7())))

	return nil
}

func (receiver *defaultRepository[TEntity]) createQuery(queryType queryType, options []FindOption) (query string, bind []any) {
	var prefix string
	switch queryType {
	case queryTypeDelete:
		prefix = "delete"
	default:
		prefix = "select"
	}

	optionsHolder := &findOptions{}
	for _, option := range options {
		option(optionsHolder)
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(prefix)
	if queryType == queryTypeSelect {
		if len(optionsHolder.selectFields) > 0 {
			queryBuilder.WriteString(" " + strings.Join(optionsHolder.selectFields, ", ") + " ")
		} else {
			queryBuilder.WriteString(" * ")
		}
	}
	queryBuilder.WriteString(" from ")
	queryBuilder.WriteString(receiver.tableName)
	if optionsHolder.tableAlias != "" {
		queryBuilder.WriteString(" " + optionsHolder.tableAlias + " ")
	}

	if len(optionsHolder.joins) > 0 {
		for table, val := range optionsHolder.joins {
			queryBuilder.WriteString(" ")
			queryBuilder.WriteString(string(val.joinType))
			queryBuilder.WriteString(" join ")
			queryBuilder.WriteString(table)
			queryBuilder.WriteString(" on ")
			queryBuilder.WriteString(val.on)
		}
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

func (receiver *defaultRepository[TEntity]) getQueryExecutor(ctx context.Context) db.QueryExecutor {
	var result db.QueryExecutor = receiver.db
	if tx := db.GetTransactionFromContext(ctx); tx != nil {
		result = tx
	}

	return result
}

func (receiver *defaultRepository[TEntity]) getQueryIssuer(ctx context.Context) db.QueryIssuer {
	var result db.QueryIssuer = receiver.db
	if tx := db.GetTransactionFromContext(ctx); tx != nil {
		result = tx
	}

	return result
}
