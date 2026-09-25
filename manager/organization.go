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

var ErrOrganizationAlreadyExists = errors.New("organization already exists")

type Organization interface {
	Create(ctx context.Context, organization *entity.Organization) error
	FindByID(ctx context.Context, organizationID uuid.UUID) (*entity.Organization, error)
	GetForUser(ctx context.Context, userID uuid.UUID) ([]*entity.Organization, error)
	Update(ctx context.Context, org *entity.Organization) error
}

func NewOrganizationManager(
	organizationRepository repo.OrganizationRepository,
) Organization {
	return &organization{
		organizationRepository: organizationRepository,
	}
}

type organization struct {
	organizationRepository repo.OrganizationRepository
}

func (receiver *organization) Create(ctx context.Context, organization *entity.Organization) error {
	err := receiver.organizationRepository.Create(ctx, organization)
	if db.IsDuplicateError(err) {
		return fmt.Errorf("%w: %w", ErrUserAlreadyExists, err)
	}
	if err != nil {
		return fmt.Errorf("failed creating organization: %w", err)
	}

	return nil
}

func (receiver *organization) FindByID(ctx context.Context, organizationID uuid.UUID) (*entity.Organization, error) {
	return receiver.organizationRepository.FindByID(ctx, organizationID)
}

func (receiver *organization) GetForUser(ctx context.Context, userID uuid.UUID) ([]*entity.Organization, error) {
	return receiver.organizationRepository.Find(
		ctx,
		repo.WithSelect([]string{"o.*"}),
		repo.WithAlias("o"),
		repo.WithJoin("organization_memberships om", "om.organization_id = o.id"),
		repo.WithWhere("om.user_id = ?", userID),
		repo.WithWhere("om.approved = true"),
	)
}

func (receiver *organization) Update(ctx context.Context, org *entity.Organization) error {
	return receiver.organizationRepository.Update(ctx, org)
}
