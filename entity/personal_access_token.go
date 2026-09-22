package entity

import (
	"time"

	"github.com/google/uuid"
)

type PersonalAccessToken struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"-"`

	Name      string `json:"name"`
	TokenHash string `json:"-"`

	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
}
