package manager

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/db/repo"
)

type AccessType string

const (
	AccessTypeNone  AccessType = "none"
	AccessTypeRead  AccessType = "read"
	AccessTypeWrite AccessType = "write"
)

type Access interface {
	CanUserAccessOrganization(ctx context.Context, userID, organizationID uuid.UUID, accessType AccessType) (bool, error)
	CanUserAccessRepository(ctx context.Context, userID, repositoryID uuid.UUID, accessType AccessType) (bool, error)
}

func NewAccessManager(
	userRepository repo.UserRepository,
) Access {
	return &access{
		userRepository: userRepository,
	}
}

type access struct {
	userRepository repo.UserRepository
}

func (receiver *access) CanUserAccessOrganization(ctx context.Context, userID, organizationID uuid.UUID, accessType AccessType) (bool, error) {
	userEntity, err := receiver.userRepository.FindByID(ctx, userID)
	if err != nil || userEntity == nil {
		return false, fmt.Errorf("failed finding user with id %s: %w", userID, err)
	}

	permission, err := receiver.userRepository.GetOrganizationPermission(ctx, userEntity, organizationID)
	if err != nil {
		return false, fmt.Errorf("failed getting org permission for user with id %s: %w", userID, err)
	}

	return receiver.decide(permission, accessType)
}

func (receiver *access) CanUserAccessRepository(ctx context.Context, userID, repositoryID uuid.UUID, accessType AccessType) (bool, error) {
	userEntity, err := receiver.userRepository.FindByID(ctx, userID)
	if err != nil || userEntity == nil {
		return false, fmt.Errorf("failed finding user with id %s: %w", userID, err)
	}

	permission, err := receiver.userRepository.GetRepositoryPermission(ctx, userEntity, repositoryID)
	if err != nil {
		return false, fmt.Errorf("failed getting repo permission for user with id %s: %w", userID, err)
	}

	return receiver.decide(permission, accessType)
}

func (receiver *access) decide(permission string, accessType AccessType) (bool, error) {
	switch permission {
	case "":
		return accessType == AccessTypeNone, nil
	case "admin":
		return accessType != AccessTypeNone, nil
	case "user":
		return accessType == AccessTypeRead, nil
	default:
		return false, fmt.Errorf("invalid permission: %s", permission)
	}
}
