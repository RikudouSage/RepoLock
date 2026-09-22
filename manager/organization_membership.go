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

var ErrOrganizationMembershipAlreadyExists = errors.New("organization membership already exists for this org and user")

type OrganizationMembership interface {
	Create(ctx context.Context, membership *entity.OrganizationMembership) error
	GetMemberships(ctx context.Context, org *entity.Organization) ([]*entity.OrganizationMembership, error)
	GetUserMembership(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*entity.OrganizationMembership, error)
	RemoveMembership(ctx context.Context, membership *entity.OrganizationMembership) error
}

func NewOrganizationMembershipManager(
	repository repo.OrganizationMembershipRepository,
) OrganizationMembership {
	return &organizationMembership{
		repository: repository,
	}
}

type organizationMembership struct {
	repository repo.OrganizationMembershipRepository
}

func (receiver *organizationMembership) Create(ctx context.Context, membership *entity.OrganizationMembership) error {
	err := receiver.repository.Create(ctx, membership)
	if db.IsDuplicateError(err) {
		return fmt.Errorf("%w: %w", ErrOrganizationMembershipAlreadyExists, err)
	}

	return err
}

func (receiver *organizationMembership) GetMemberships(ctx context.Context, org *entity.Organization) ([]*entity.OrganizationMembership, error) {
	return receiver.repository.Find(
		ctx,
		repo.WithWhere("organization_id = ?", org.ID),
		repo.WithOrderBy("permission = 'admin'", "asc"),
	)
}

func (receiver *organizationMembership) GetUserMembership(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*entity.OrganizationMembership, error) {
	items, err := receiver.repository.Find(
		ctx,
		repo.WithWhere("organization_id = ?", orgID),
		repo.WithWhere("user_id = ?", userID),
	)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	return items[0], nil
}

func (receiver *organizationMembership) RemoveMembership(ctx context.Context, membership *entity.OrganizationMembership) error {
	return receiver.repository.Delete(ctx, repo.WithWhere("id = ?", membership.ID))
}
