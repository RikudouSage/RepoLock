package manager

import (
	"context"
	"errors"
	"fmt"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/entity"
)

var ErrOrganizationMembershipAlreadyExists = errors.New("organization membership already exists for this org and user")

type OrganizationMembership interface {
	Create(ctx context.Context, membership *entity.OrganizationMembership) error
	GetMemberships(ctx context.Context, org *entity.Organization) ([]*entity.OrganizationMembership, error)
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
