package dto

import (
	"errors"

	"github.com/google/uuid"
)

type CreateRepositoryRequest struct {
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

func (receiver *CreateRepositoryRequest) Validate() error {
	if receiver.URL == "" {
		return errors.New("url is required")
	}
	if receiver.OrganizationID == uuid.Nil {
		return errors.New("organization_id is required")
	}

	return nil
}
