package dto

import (
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/entity"
)

type OrganizationPermission struct {
	OrganizationID uuid.UUID         `json:"organization_id"`
	UserID         uuid.UUID         `json:"user_id"`
	Permission     entity.Permission `json:"permission"`
	Approved       bool              `json:"approved"`
}

type RepositoryPermission struct {
	UserID           uuid.UUID         `json:"user_id"`
	RepositoryID     uuid.UUID         `json:"repository_id"`
	Permission       entity.Permission `json:"permission"`
	FromOrganization bool              `json:"from_organization"`
}
