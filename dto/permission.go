package dto

import (
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/entity"
)

type PermissionType string

const (
	PermissionTypeGroup PermissionType = "group"
	PermissionTypeUser  PermissionType = "user"
)

type Permission struct {
	Type PermissionType `json:"type"`

	Permission entity.Permission `json:"permission"`
	UserID     uuid.UUID         `json:"user_id"`

	OrganizationID uuid.UUID `json:"organization_id,omitzero"`
	RepositoryID   uuid.UUID `json:"repository_id,omitzero"`
	Approved       bool      `json:"approved"`
}
