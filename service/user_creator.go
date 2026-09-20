package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/entity"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

type UserCreator interface {
	Create(ctx context.Context, userRegistration dto.UserRegister) (*entity.User, error)
}

func NewUserCreator(
	userManager manager.User,
	organizationManager manager.Organization,
	organizationMembershipRepository repo.OrganizationMembershipRepository,
	transactionCreator db.TransactionCreator,
	passwordHasher PasswordHasher,
) UserCreator {
	return &userCreator{
		userManager:                      userManager,
		organizationManager:              organizationManager,
		organizationMembershipRepository: organizationMembershipRepository,
		transactionCreator:               transactionCreator,
		passwordHasher:                   passwordHasher,
	}
}

type userCreator struct {
	userManager                      manager.User
	organizationManager              manager.Organization
	organizationMembershipRepository repo.OrganizationMembershipRepository
	transactionCreator               db.TransactionCreator
	passwordHasher                   PasswordHasher
}

func (receiver *userCreator) Create(ctx context.Context, userRegistration dto.UserRegister) (*entity.User, error) {
	tx, err := receiver.transactionCreator.Create()
	if err != nil {
		return nil, fmt.Errorf("failed creating transaction: %w", err)
	}

	ctx = db.WithTransaction(ctx, tx)

	hash, err := receiver.passwordHasher.Hash(userRegistration.Password)
	if err != nil {
		return nil, fmt.Errorf("failed hashing password: %w", err)
	}

	user := &entity.User{
		Username:     userRegistration.Username,
		PasswordHash: &hash,
		Name:         userRegistration.Name,
	}
	if err = receiver.userManager.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	permission := entity.PermissionUser
	approved := false

	var org *entity.Organization
	if userRegistration.OrganizationID != uuid.Nil {
		org, err = receiver.organizationManager.FindByID(ctx, userRegistration.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("failed finding organization: %w", err)
		}
	} else {
		permission = entity.PermissionAdmin
		approved = true
		org = &entity.Organization{
			Name: fmt.Sprintf("Organization of user %s", user.Name),
		}
		err = receiver.organizationManager.Create(ctx, org)
		if err != nil {
			return nil, fmt.Errorf("failed creating organization for new user: %w", err)
		}
	}

	membership := &entity.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Permission:     permission,
		Approved:       approved,
	}
	if err := receiver.organizationMembershipRepository.Create(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed creating organization membership: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed committing transaction: %w", err)
	}

	return user, nil
}
