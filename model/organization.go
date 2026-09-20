package model

import "github.com/google/uuid"

type Organization struct {
	ID uuid.UUID

	Name string
}
