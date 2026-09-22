package manager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/entity"
	"go.chrastecky.dev/repolock/service"
)

type PATConfig struct {
	Name      string
	ExpiresAt *time.Time
}

type PersonalAccessToken interface {
	CreateForUser(ctx context.Context, user *entity.User, opts *PATConfig) (*entity.PersonalAccessToken, string, error)
	FindByToken(ctx context.Context, authorization string) (*entity.PersonalAccessToken, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.PersonalAccessToken, error)
	DeleteToken(ctx context.Context, pat *entity.PersonalAccessToken) error
	FindAllForUser(ctx context.Context, user *entity.User) ([]*entity.PersonalAccessToken, error)
	UpdateLastUsedAt(ctx context.Context, pat *entity.PersonalAccessToken, dateTime time.Time) error
}

func NewPersonalAccessTokenManager(
	repository repo.PersonalAccessTokenRepository,
	passwordHasher service.PasswordHasher,
	passwordVerifier service.PasswordVerifier,
	randomStringGenerator service.RandomStringGenerator,
) PersonalAccessToken {
	return &personalAccessToken{
		repository:            repository,
		passwordHasher:        passwordHasher,
		passwordVerifier:      passwordVerifier,
		randomStringGenerator: randomStringGenerator,
		now:                   time.Now,
	}
}

type personalAccessToken struct {
	repository            repo.PersonalAccessTokenRepository
	passwordHasher        service.PasswordHasher
	passwordVerifier      service.PasswordVerifier
	randomStringGenerator service.RandomStringGenerator
	now                   func() time.Time
}

func (receiver *personalAccessToken) CreateForUser(ctx context.Context, user *entity.User, opts *PATConfig) (*entity.PersonalAccessToken, string, error) {
	rawToken := receiver.randomStringGenerator.GenerateRandomString(64)
	hash, err := receiver.passwordHasher.Hash(rawToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed hashing password: %w", err)
	}

	pat := &entity.PersonalAccessToken{
		UserID:    user.ID,
		Name:      opts.Name,
		TokenHash: hash,
		CreatedAt: receiver.now(),
		ExpiresAt: opts.ExpiresAt,
	}

	if err := receiver.repository.Create(ctx, pat); err != nil {
		return nil, "", fmt.Errorf("failed persisting personal access token: %w", err)
	}

	return pat, fmt.Sprintf("pat_%s_%s", pat.ID, rawToken), nil
}

func (receiver *personalAccessToken) FindByToken(ctx context.Context, token string) (*entity.PersonalAccessToken, error) {
	parts := strings.SplitN(token, "_", 3)

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format, expected 3 parts split by _")
	}

	if parts[0] != "pat" {
		return nil, fmt.Errorf("invalid token format, expected pat, got '%s'", parts[0])
	}

	tokenID, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token format, expected uuid, got '%s'", parts[1])
	}

	pat, err := receiver.repository.FindByID(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed persisting personal access token: %w", err)
	}

	if pat == nil {
		return nil, fmt.Errorf("personal access token not found")
	}

	verify, err := receiver.passwordVerifier.Verify(parts[2], pat.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed verifying personal access token: %w", err)
	}

	if !verify {
		return nil, fmt.Errorf("invalid personal access token")
	}

	return pat, nil
}

func (receiver *personalAccessToken) FindByID(ctx context.Context, id uuid.UUID) (*entity.PersonalAccessToken, error) {
	return receiver.repository.FindByID(ctx, id)
}

func (receiver *personalAccessToken) DeleteToken(ctx context.Context, pat *entity.PersonalAccessToken) error {
	return receiver.repository.Delete(ctx, repo.WithWhere("id = ?", pat.ID))
}

func (receiver *personalAccessToken) FindAllForUser(ctx context.Context, user *entity.User) ([]*entity.PersonalAccessToken, error) {
	return receiver.repository.Find(
		ctx,
		repo.WithWhere("user_id = ?", user.ID),
		repo.WithOrderBy("name", "asc"),
	)
}

func (receiver *personalAccessToken) UpdateLastUsedAt(ctx context.Context, pat *entity.PersonalAccessToken, dateTime time.Time) error {
	return receiver.repository.UpdateLastUsedAt(ctx, pat, dateTime)
}
