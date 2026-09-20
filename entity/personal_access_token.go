package entity

import (
	"time"

	"github.com/google/uuid"
)

type PersonalAccessToken struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Name      string
	TokenHash string

	CreatedAt  time.Time
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
}
