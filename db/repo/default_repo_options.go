package repo

type findOptions struct {
	selectFields     []string
	where            []string
	limit            uint
	orderByField     string
	orderByDirection string
	bindValues       []any
	joins            map[string]string
	tableAlias       string
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

func WithJoin(table, on string) FindOption {
	return func(options *findOptions) {
		if options.joins == nil {
			options.joins = make(map[string]string)
		}
		options.joins[table] = on
	}
}

func WithSelect(fields []string) FindOption {
	return func(options *findOptions) {
		options.selectFields = fields
	}
}

func WithAlias(alias string) FindOption {
	return func(options *findOptions) {
		options.tableAlias = alias
	}
}
