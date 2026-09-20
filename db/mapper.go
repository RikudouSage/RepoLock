package db

import (
	"fmt"
	"reflect"
	"slices"

	strcase "github.com/iancoleman/strcase"
)

type Scannable interface {
	Scan(dest ...any) error
	Columns() ([]string, error)
}

type Mapper interface {
	Map(rows Scannable, result any) error
}

func NewMapper() Mapper {
	return &mapper{}
}

type mapper struct {
}

func (receiver *mapper) Map(scannable Scannable, result any) error {
	ref := reflect.ValueOf(result)
	if ref.Kind() != reflect.Pointer {
		return fmt.Errorf("result must be a pointer")
	}
	if ref.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("result must be a struct, %T given", result)
	}

	mappedFieldNames := make([]string, ref.Elem().NumField())
	for i := range ref.Elem().NumField() {
		mappedFieldNames[i] = strcase.ToSnake(ref.Elem().Type().Field(i).Name)
	}

	columns, err := scannable.Columns()
	if err != nil {
		return fmt.Errorf("failed getting column list: %w", err)
	}

	binders := make([]any, len(columns))
	for i, column := range columns {
		idx := slices.Index(mappedFieldNames, column)
		if idx == -1 {
			return fmt.Errorf("column %s not found in %T", column, result)
		}

		binders[i] = ref.Elem().Field(idx).Addr().Interface()
	}

	err = scannable.Scan(binders...)
	if err != nil {
		return fmt.Errorf("failed scanning row into %T: %w", result, err)
	}

	return nil
}
