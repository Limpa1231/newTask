package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LegalEntity struct {
	UUID           uuid.UUID `json:"uuid"`
	Name           string    `json:"name"`
	FederationUUID uuid.UUID `json:"federation_uuid"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LegalEntityUpdate struct {
	Name      *string `json:"name,omitempty"`
	UpdatedBy string  `json:"updated_by"`
}

