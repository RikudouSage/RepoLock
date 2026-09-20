package model

import (
	"time"

	"github.com/google/uuid"
)

type Lock struct {
	ID           uuid.UUID
	RepositoryID uuid.UUID
	UserID       uuid.UUID

	Pattern string
	Reason  string

	CreatedAt time.Time
	ExpiresAt *time.Time
}
