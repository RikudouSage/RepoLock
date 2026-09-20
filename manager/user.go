package manager

import (
	"context"
	"errors"
	"fmt"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/entity"
	"go.chrastecky.dev/repolock/service"
)

const invalidPasswordToCheck = "$argon2id$v=19$m=262144,t=3,p=4$xTYebMoCMUx411QXUzbskQ$T7Q9rO7n5CikSQG6OZkcwe+P0sqxpV2GB4S/NJAunVA"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrPasswordNotEnabledForUser = errors.New("password not enabled for user")
var ErrInvalidPassword = errors.New("invalid password")

type User interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsernameAndPassword(ctx context.Context, username string, password string) (*entity.User, error)
}

func NewUserManager(
	userRepository repo.UserRepository,
	passwordVerifier service.PasswordVerifier,
) User {
	return &user{
		userRepository:   userRepository,
		passwordVerifier: passwordVerifier,
	}
}

type user struct {
	userRepository   repo.UserRepository
	passwordVerifier service.PasswordVerifier
}

func (receiver *user) Create(ctx context.Context, user *entity.User) error {
	err := receiver.userRepository.Create(ctx, user)
	if db.IsDuplicateError(err) {
		return fmt.Errorf("%w: %w", ErrUserAlreadyExists, err)
	}
	if err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}

	return nil
}

func (receiver *user) FindByUsernameAndPassword(ctx context.Context, username string, password string) (*entity.User, error) {
	userEntities, err := receiver.userRepository.Find(
		ctx,
		repo.WithWhere("username = ?", username),
		repo.WithLimit(1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed finding user by username: %w", err)
	}

	if len(userEntities) == 0 {
		_, _ = receiver.passwordVerifier.Verify("invalid password to test", invalidPasswordToCheck)
		return nil, ErrUserNotFound
	}
	if userEntities[0].PasswordHash == nil {
		_, _ = receiver.passwordVerifier.Verify("invalid password to test", invalidPasswordToCheck)
		return nil, ErrPasswordNotEnabledForUser
	}

	valid, err := receiver.passwordVerifier.Verify(password, *userEntities[0].PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed verifying password: %w", err)
	}

	if !valid {
		return nil, ErrInvalidPassword
	}

	return userEntities[0], nil
}
