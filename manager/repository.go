package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/entity"
	"go.chrastecky.dev/repolock/service"
)

var ErrRepositoryNotFound = errors.New("repository not found")
var ErrRepositoryAlreadyExists = errors.New("repository already exists")

type Repository interface {
	GetReposForUser(ctx context.Context, user *entity.User) ([]*entity.Repository, error)
	GetRepoForUser(ctx context.Context, user *entity.User, repoID uuid.UUID) (*entity.Repository, error)
	CreateRepository(ctx context.Context, id uuid.UUID, url string, name string) (*entity.Repository, error)
}

func NewRepositoryManager(
	repo repo.RepositoryRepository,
	urlNormalizer service.RepositoryURLNormalizer,
) Repository {
	return &repository{
		repository:    repo,
		urlNormalizer: urlNormalizer,
	}
}

type repository struct {
	repository    repo.RepositoryRepository
	urlNormalizer service.RepositoryURLNormalizer
}

func (receiver *repository) GetReposForUser(ctx context.Context, user *entity.User) ([]*entity.Repository, error) {
	return receiver.repository.FindForUser(
		ctx,
		user,
		repo.WithOrderBy("identifier", "asc"),
	)
}

func (receiver *repository) GetRepoForUser(ctx context.Context, user *entity.User, repoID uuid.UUID) (*entity.Repository, error) {
	repos, err := receiver.repository.FindForUser(
		ctx,
		user,
		repo.WithWhere("r.id = ?", repoID),
		repo.WithLimit(1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed fetching repositories: %w", err)
	}

	if len(repos) == 0 {
		return nil, ErrRepositoryNotFound
	}

	return repos[0], nil
}

func (receiver *repository) CreateRepository(ctx context.Context, organizationID uuid.UUID, url string, name string) (*entity.Repository, error) {
	identifier, err := receiver.urlNormalizer.Normalize(url)
	if err != nil {
		return nil, fmt.Errorf("failed normalizing url: %w", err)
	}
	if name == "" {
		name = identifier
	}

	repoEntity := &entity.Repository{
		OrganizationID: organizationID,
		Name:           name,
		Identifier:     identifier,
	}

	if err = receiver.repository.Create(ctx, repoEntity); err != nil {
		if db.IsDuplicateError(err) {
			return nil, fmt.Errorf("%w: %w", ErrRepositoryAlreadyExists, err)
		}

		return nil, fmt.Errorf("failed creating repository: %w", err)
	}

	return repoEntity, nil
}
