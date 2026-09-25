package repo

import (
	"database/sql"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type VCSIdentityRepository interface {
	DefaultRepository[entity.VCSIdentity]
}

func NewVCSIdentityRepository(
	db *sql.DB,
	formatter db.QueryFormatter,
	mapper db.Mapper,
) VCSIdentityRepository {
	repo := &vcsIdentityRepository{}
	repo.defaultRepository = newDefaultRepository[entity.VCSIdentity](
		db,
		"vcs_identities",
		formatter,
		mapper,
	)

	return repo
}

type vcsIdentityRepository struct {
	*defaultRepository[entity.VCSIdentity]
}
