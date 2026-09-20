package manager

import (
	"context"
	"errors"
	"fmt"

	"go.chrastecky.dev/repolock/config/data"
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/entity"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type User interface {
	Create(ctx context.Context, user *entity.User) error
}

func NewUserManager(
	userRepository repo.UserRepository,
) User {
	return &user{
		userRepository: userRepository,
	}
}

type user struct {
	dbType         data.DatabaseType
	userRepository repo.UserRepository
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
