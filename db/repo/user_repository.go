package repo

import (
	"database/sql"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type UserRepository interface {
	DefaultRepository[entity.User]
}

func NewUserRepository(
	db *sql.DB,
	formatter db.QueryFormatter,
	mapper db.Mapper,
) UserRepository {
	repo := &userRepository{}
	repo.defaultRepository = newDefaultRepository[entity.User](
		db,
		"users",
		formatter,
		mapper,
	)

	return repo
}

type userRepository struct {
	*defaultRepository[entity.User]
}
