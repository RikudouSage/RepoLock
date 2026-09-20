package entity

import "github.com/google/uuid"

type Repository struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID

	Name       string
	Identifier string
}
