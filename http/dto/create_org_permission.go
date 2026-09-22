package dto

import (
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/entity"
)

type CreateOrganizationPermission struct {
	UserID     uuid.UUID         `json:"user_id"`
	Permission entity.Permission `json:"permission"`
	Approved   bool              `json:"approved"`
}
