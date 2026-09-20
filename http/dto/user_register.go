package dto

import "github.com/google/uuid"

type UserRegister struct {
	Name           string    `json:"name"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	OrganizationID uuid.UUID `json:"organization_id"`
}
