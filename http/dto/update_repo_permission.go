package dto

import "go.chrastecky.dev/repolock/entity"

type UpdateRepositoryPermission struct {
	Permission entity.Permission `json:"permission"`
}
