package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type UserRepository interface {
	DefaultRepository[entity.User]

	GetOrganizationPermission(ctx context.Context, user *entity.User, organizationID uuid.UUID) (string, error)
	GetRepositoryPermission(ctx context.Context, user *entity.User, repositoryID uuid.UUID) (string, error)
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

func (receiver *userRepository) GetOrganizationPermission(ctx context.Context, user *entity.User, organizationID uuid.UUID) (string, error) {
	database := receiver.getQueryIssuer(ctx)

	query := "select permission from organization_memberships where organization_id = ? and user_id = ? and approved = true"
	query = receiver.queryFormatter.FormatQuery(query)

	var permission string
	err := database.QueryRowContext(ctx, query, organizationID, user.ID).Scan(&permission)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed getting org permission: %w", err)
	}

	return permission, nil
}

func (receiver *userRepository) GetRepositoryPermission(ctx context.Context, user *entity.User, repositoryID uuid.UUID) (string, error) {
	database := receiver.getQueryIssuer(ctx)
	query := `select permission
				from (
						 select
							 rp.permission,
							 case rp.permission when 'admin' then 2 else 1 end as priority
						 from repository_permissions rp
						 where rp.repository_id = ?
						   and rp.user_id = ?
				
						 union all
				
						 select
							 om.permission,
							 case om.permission when 'admin' then 2 else 1 end as priority
						 from repositories r
								  join organization_memberships om
									   on om.organization_id = r.organization_id
						 where r.id = ?
						   and om.user_id = ?
						   and om.approved = true
					 ) permissions
				order by priority desc
				limit 1`
	query = receiver.queryFormatter.FormatQuery(query)
	var permission string
	err := database.QueryRowContext(ctx, query, repositoryID, user.ID, repositoryID, user.ID).Scan(&permission)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed getting org permission: %w", err)
	}

	return permission, nil
}
