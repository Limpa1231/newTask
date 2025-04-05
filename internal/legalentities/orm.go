package legalentities

import (
	"time"

	"github.com/google/uuid"
)

type LegalEntity struct {
	UUID      uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string     `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time  `gorm:"type:timestamptz;default:now();not null"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;default:now();not null"`
	DeletedAt *time.Time `gorm:"type:timestamptz;index;default:NULL"`
}

func (LegalEntity) TableName() string {
	return "legal_entities"
}
