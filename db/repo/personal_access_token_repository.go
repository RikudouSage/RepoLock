package repo

import (
	"database/sql"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type PersonalAccessTokenRepository interface {
	DefaultRepository[entity.PersonalAccessToken]
}

func NewPersonalAccessTokenRepository(
	db *sql.DB,
	formatter db.QueryFormatter,
	mapper db.Mapper,
) PersonalAccessTokenRepository {
	repo := &personalAccessTokenRepository{}
	repo.defaultRepository = newDefaultRepository[entity.PersonalAccessToken](
		db,
		"personal_access_tokens",
		formatter,
		mapper,
	)

	return repo
}

type personalAccessTokenRepository struct {
	*defaultRepository[entity.PersonalAccessToken]
}
