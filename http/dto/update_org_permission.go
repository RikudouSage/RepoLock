package dto

import "go.chrastecky.dev/repolock/entity"

type UpdateOrganizationPermission struct {
	Permission entity.Permission `json:"permission"`
	Approved   bool              `json:"approved"`
}
