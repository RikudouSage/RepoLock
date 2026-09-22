package repo

import (
	"context"
	"database/sql"
	"fmt"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type RepositoryRepository interface {
	DefaultRepository[entity.Repository]

	FindForUser(ctx context.Context, user *entity.User, options ...FindOption) ([]*entity.Repository, error)
}

func NewRepositoryRepository(
	db *sql.DB,
	formatter db.QueryFormatter,
	mapper db.Mapper,
) RepositoryRepository {
	repo := &repositoryRepository{}
	repo.defaultRepository = newDefaultRepository[entity.Repository](
		db,
		"repositories",
		formatter,
		mapper,
	)

	return repo
}

type repositoryRepository struct {
	*defaultRepository[entity.Repository]
}

func (receiver *repositoryRepository) FindForUser(ctx context.Context, user *entity.User, options ...FindOption) ([]*entity.Repository, error) {
	database := receiver.getQueryIssuer(ctx)

	options = append([]FindOption{
		WithWhere("om.user_id = ?", user.ID),
		WithSkipAutoTableName(true),
	}, options...)

	query := `select r.*
				from repositories r
				join organization_memberships om
					on om.organization_id = r.organization_id`
	query, bind := receiver.createQuery(query, options)
	query = receiver.queryFormatter.FormatQuery(query)

	rows, err := database.QueryContext(ctx, query, bind...)
	if err != nil {
		return nil, fmt.Errorf("failed getting repositories for user %s: %w", user.ID, err)
	}
	defer rows.Close()

	result := make([]*entity.Repository, 0)
	for rows.Next() {
		var item entity.Repository
		if err = receiver.mapper.MapOntoStruct(rows, &item); err != nil {
			return nil, fmt.Errorf("failed mapping repository for user %s: %w", user.ID, err)
		}
		result = append(result, &item)
	}

	return result, nil
}
