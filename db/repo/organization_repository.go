package repo

import (
	"database/sql"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type OrganizationRepository interface {
	DefaultRepository[entity.Organization]
}

func NewOrganizationRepository(
	db *sql.DB,
	formatter db.QueryFormatter,
	mapper db.Mapper,
) OrganizationRepository {
	repo := &organizationRepository{}
	repo.defaultRepository = newDefaultRepository[entity.Organization](
		db,
		"organizations",
		formatter,
		mapper,
	)

	return repo
}

type organizationRepository struct {
	*defaultRepository[entity.Organization]
}
