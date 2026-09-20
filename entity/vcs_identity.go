package entity

import "github.com/google/uuid"

type VCSIdentity struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Identity string
}
