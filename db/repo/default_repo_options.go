package repo

type findOptions struct {
	where             []string
	limit             uint
	orderByField      string
	orderByDirection  string
	bindValues        []any
	skipAutoTableName bool
}

type FindOption func(*findOptions)

func WithWhere(where string, bind ...any) FindOption {
	return func(options *findOptions) {
		options.where = append(options.where, where)
		options.bindValues = append(options.bindValues, bind...)
	}
}

func WithLimit[TInt ~int | ~uint](limit TInt) FindOption {
	return func(options *findOptions) {
		options.limit = uint(limit)
	}
}

func WithOrderBy(column, direction string) FindOption {
	return func(options *findOptions) {
		options.orderByField = column
		options.orderByDirection = direction
	}
}

func WithSkipAutoTableName(skip bool) FindOption {
	return func(options *findOptions) {
		options.skipAutoTableName = skip
	}
}
