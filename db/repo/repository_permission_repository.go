package repo

import (
	"database/sql"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type RepositoryPermissionRepository interface {
	DefaultRepository[entity.RepositoryPermission]
}

func NewRepositoryPermissionRepository(
	db *sql.DB,
	formatter db.QueryFormatter,
	mapper db.Mapper,
) RepositoryPermissionRepository {
	repo := &repositoryPermissionRepository{}
	repo.defaultRepository = newDefaultRepository[entity.RepositoryPermission](
		db,
		"repository_permissions",
		formatter,
		mapper,
	)

	return repo
}

type repositoryPermissionRepository struct {
	*defaultRepository[entity.RepositoryPermission]
}
