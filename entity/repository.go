package entity

import "github.com/google/uuid"

type Repository struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"-"`

	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}
