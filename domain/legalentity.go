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

type LegalEntitiesRepository interface {
	Create(ctx context.Context, entity *LegalEntity) error
	GetByFederation(ctx context.Context, federationID uuid.UUID) ([]*LegalEntity, error)
	Update(ctx context.Context, id uuid.UUID, name string, updatedBy string) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

type LegalEntitiesService interface {
	Create(ctx context.Context, entity *LegalEntity) error
	GetByFederation(ctx context.Context, federationID uuid.UUID) ([]*LegalEntity, error)
	Update(ctx context.Context, id uuid.UUID, name string, updatedBy string) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}
