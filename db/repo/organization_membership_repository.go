package repo

import (
	"database/sql"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type OrganizationMembershipRepository interface {
	DefaultRepository[entity.OrganizationMembership]
}

func NewOrganizationMembershipRepository(
	db *sql.DB,
	formatter db.QueryFormatter,
	mapper db.Mapper,
) OrganizationMembershipRepository {
	repo := &organizationMembershipRepository{}
	repo.defaultRepository = newDefaultRepository[entity.OrganizationMembership](
		db,
		"organization_memberships",
		formatter,
		mapper,
	)

	return repo
}

type organizationMembershipRepository struct {
	*defaultRepository[entity.OrganizationMembership]
}
