package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/entity"
	"go.chrastecky.dev/repolock/http/dto"
)

var ErrPermissionAlreadyExists = errors.New("permission for this user already exists")

type RepositoryPermission interface {
	GetPermissions(ctx context.Context, repository *entity.Repository) ([]*dto.RepositoryPermission, error)
	Create(ctx context.Context, perm *entity.RepositoryPermission) error
	GetUserPermission(ctx context.Context, repoID uuid.UUID, userID uuid.UUID) (*entity.RepositoryPermission, error)
	Delete(ctx context.Context, perm *entity.RepositoryPermission) error
	Update(ctx context.Context, perm *entity.RepositoryPermission) error
}

func NewRepositoryPermissionManager(
	repoPermissionRepository repo.RepositoryPermissionRepository,
	orgMembershipRepository repo.OrganizationMembershipRepository,
) RepositoryPermission {
	return &repositoryPermission{
		repositoryPermissionRepository:   repoPermissionRepository,
		organizationMembershipRepository: orgMembershipRepository,
	}
}

type repositoryPermission struct {
	repositoryPermissionRepository   repo.RepositoryPermissionRepository
	organizationMembershipRepository repo.OrganizationMembershipRepository
}

func (receiver *repositoryPermission) GetPermissions(ctx context.Context, repository *entity.Repository) ([]*dto.RepositoryPermission, error) {
	repoPermissions, err := receiver.repositoryPermissionRepository.Find(
		ctx,
		repo.WithWhere("repository_id = ?", repository.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("find repository permissions: %w", err)
	}

	orgMemberships, err := receiver.organizationMembershipRepository.Find(
		ctx,
		repo.WithAlias("om"),
		repo.WithSelect([]string{"om.*"}),
		repo.WithJoin("repositories r", "om.organization_id = r.organization_id"),
		repo.WithWhere("r.id = ?", repository.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("find organization permissions: %w", err)
	}

	return lo.Concat(
		lo.Map(repoPermissions, func(item *entity.RepositoryPermission, _ int) *dto.RepositoryPermission {
			return &dto.RepositoryPermission{
				UserID:           item.UserID,
				RepositoryID:     item.RepositoryID,
				Permission:       item.Permission,
				FromOrganization: false,
			}
		}),
		lo.Map(orgMemberships, func(item *entity.OrganizationMembership, _ int) *dto.RepositoryPermission {
			return &dto.RepositoryPermission{
				UserID:           item.UserID,
				RepositoryID:     repository.ID,
				Permission:       item.Permission,
				FromOrganization: true,
			}
		}),
	), nil
}

func (receiver *repositoryPermission) Create(ctx context.Context, perm *entity.RepositoryPermission) error {
	err := receiver.repositoryPermissionRepository.Create(ctx, perm)
	if db.IsDuplicateError(err) {
		return fmt.Errorf("%w: %w", ErrPermissionAlreadyExists, err)
	}

	return err
}

func (receiver *repositoryPermission) GetUserPermission(ctx context.Context, repoID uuid.UUID, userID uuid.UUID) (*entity.RepositoryPermission, error) {
	items, err := receiver.repositoryPermissionRepository.Find(
		ctx,
		repo.WithWhere("repository_id = ?", repoID),
		repo.WithWhere("user_id = ?", userID),
		repo.WithLimit(1),
	)
	if err != nil {
		return nil, fmt.Errorf("find repository permissions: %w", err)
	}

	if len(items) == 0 {
		return nil, nil
	}

	return items[0], nil
}

func (receiver *repositoryPermission) Delete(ctx context.Context, perm *entity.RepositoryPermission) error {
	return receiver.repositoryPermissionRepository.Delete(ctx, repo.WithWhere("id = ?", perm.ID))
}

func (receiver *repositoryPermission) Update(ctx context.Context, perm *entity.RepositoryPermission) error {
	return receiver.repositoryPermissionRepository.Update(ctx, perm)
}
