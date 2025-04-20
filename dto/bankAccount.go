package dto

import (
	"time"

	"github.com/google/uuid"
)

type BankAccountDTO struct {
	UUID            uuid.UUID  `json:"uuid"`
	LegalEntityUUID *uuid.UUID `json:"legal_entity_uuid,omitempty"`
	BIC             string     `json:"bic"`
	BankName        string     `json:"bank"`
	BankAddress     string     `json:"address"`
	CorrAccount     string     `json:"corr_account"`
	PaymentAccount  string     `json:"current_account"`
	Currency        string     `json:"currency"`
	Comment         string     `json:"comment"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
