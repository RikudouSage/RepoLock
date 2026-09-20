package entity

import "github.com/google/uuid"

type OrganizationMembership struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID

	Permission Permission
	Approved   bool
}
