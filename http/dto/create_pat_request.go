package dto

import "time"

type CreatePersonalAccessTokenRequest struct {
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expires_at"`
}
