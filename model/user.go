package model

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID

	Username     string
	PasswordHash *string

	Name string
}
