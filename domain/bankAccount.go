package domain

import (
	"time"

	"github.com/google/uuid"
)

type BankAccount struct {
	UUID            uuid.UUID `validate:"uuid"`
	LegalEntityUUID uuid.UUID `validate:"uuid"`
	BIC             string    `validate:"len=9"`
	BankName        string    `validate:"required,max=255"`
	BankAddress     string    `validate:"max=500"`
	CorrAccount     string    `validate:"len=20"`
	PaymentAccount  string    `validate:"len=20"`
	Currency        string    `validate:"len=3"`
	Comment         string    `validate:"max=1000"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
