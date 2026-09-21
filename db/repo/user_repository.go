package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type UserRepository interface {
	DefaultRepository[entity.User]

	GetOrganizationPermission(ctx context.Context, user *entity.User, organizationID uuid.UUID) (string, error)
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

	return permission, nil
}
