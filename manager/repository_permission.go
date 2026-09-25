package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/dto"
	"go.chrastecky.dev/repolock/entity"
)

var ErrPermissionAlreadyExists = errors.New("permission for this user already exists")

type RepositoryPermission interface {
	GetPermissions(ctx context.Context, repository *entity.Repository) ([]*dto.Permission, error)
	Create(ctx context.Context, perm *entity.RepositoryPermission) error
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

func (receiver *repositoryPermission) GetPermissions(ctx context.Context, repository *entity.Repository) ([]*dto.Permission, error) {
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
		lo.Map(repoPermissions, func(item *entity.RepositoryPermission, _ int) *dto.Permission {
			return &dto.Permission{
				Type:         dto.PermissionTypeUser,
				Permission:   item.Permission,
				UserID:       item.UserID,
				RepositoryID: item.RepositoryID,
				Approved:     true,
			}
		}),
		lo.Map(orgMemberships, func(item *entity.OrganizationMembership, _ int) *dto.Permission {
			return &dto.Permission{
				Type:           dto.PermissionTypeGroup,
				Permission:     item.Permission,
				UserID:         item.UserID,
				OrganizationID: item.OrganizationID,
				Approved:       item.Approved,
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
