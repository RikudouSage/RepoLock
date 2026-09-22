package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
)

type PersonalAccessTokenRepository interface {
	DefaultRepository[entity.PersonalAccessToken]

	UpdateLastUsedAt(ctx context.Context, pat *entity.PersonalAccessToken, dateTime time.Time) error
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

func (receiver *personalAccessTokenRepository) UpdateLastUsedAt(ctx context.Context, pat *entity.PersonalAccessToken, dateTime time.Time) error {
	database := receiver.getQueryExecutor(ctx)
	_, err := database.ExecContext(
		ctx,
		"update personal_access_tokens set last_used_at = ? where id = ?",
		dateTime,
		pat.ID,
	)

	if err != nil {
		return fmt.Errorf("failed updating last used at: %w", err)
	}

	return nil
}
