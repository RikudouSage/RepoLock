package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/entity"
)

var ErrIdentityAlreadyExists = errors.New("the specified identity already exists")

type VCSIdentity interface {
	GetForUser(ctx context.Context, id uuid.UUID) ([]*entity.VCSIdentity, error)
	CreateForUser(ctx context.Context, identity string, userID uuid.UUID) (*entity.VCSIdentity, error)
}

func NewVCSIdentityManager(
	vcsRepo repo.VCSIdentityRepository,
) VCSIdentity {
	return &vcsIdentity{
		vcsRepo: vcsRepo,
	}
}

type vcsIdentity struct {
	vcsRepo repo.VCSIdentityRepository
}

func (receiver *vcsIdentity) GetForUser(ctx context.Context, id uuid.UUID) ([]*entity.VCSIdentity, error) {
	return receiver.vcsRepo.Find(
		ctx,
		repo.WithWhere("user_id = ?", id),
	)
}

func (receiver *vcsIdentity) CreateForUser(ctx context.Context, identity string, userID uuid.UUID) (*entity.VCSIdentity, error) {
	item := &entity.VCSIdentity{
		UserID:   userID,
		Identity: identity,
	}
	err := receiver.vcsRepo.Create(ctx, item)
	if err != nil {
		if db.IsDuplicateError(err) {
			return nil, fmt.Errorf("%w: %w", ErrIdentityAlreadyExists, err)
		}

		return nil, fmt.Errorf("failed to create vcs identity: %w", err)
	}

	return item, nil
}
