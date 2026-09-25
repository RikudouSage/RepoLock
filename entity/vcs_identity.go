package entity

import "github.com/google/uuid"

type VCSIdentity struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`

	Identity string `json:"identity"`
}
