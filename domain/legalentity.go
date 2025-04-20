package domain

import (
	"time"

	"github.com/google/uuid"
)

type LegalEntity struct {
	UUID         uuid.UUID     `json:"uuid"`
	Name         string        `json:"name"`
	BankAccounts []BankAccount `json:"bank_accounts" gorm:"foreignKey:LegalEntityUUID;references:UUID"`
	CreatedBy    string        `json:"created_by"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type LegalEntityUpdate struct {
	Name      *string `json:"name,omitempty"`
	UpdatedBy string  `json:"updated_by"`
}
