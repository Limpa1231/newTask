package dto

import (
	"time"

	"github.com/google/uuid"
)

type LegalEntitiesDTO struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
