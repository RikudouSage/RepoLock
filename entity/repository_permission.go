package entity

import "github.com/google/uuid"

type RepositoryPermission struct {
	ID           uuid.UUID
	RepositoryID uuid.UUID
	UserID       uuid.UUID

	Permission Permission
}
