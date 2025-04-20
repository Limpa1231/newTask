package legalentities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LegalEntity struct {
	UUID         uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string        `gorm:"type:varchar(50);not null"`
	BankAccounts []BankAccount `gorm:"foreignKey:LegalEntityUUID;references:UUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt    time.Time     `gorm:"type:timestamptz;default:now();not null"`
	UpdatedAt    time.Time     `gorm:"type:timestamptz;default:now();not null"`
	DeletedAt    *time.Time    `gorm:"type:timestamptz;index;default:NULL"`
}

func (LegalEntity) TableName() string {
	return "legal_entities"
}

type BankAccount struct {
	UUID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	LegalEntityUUID uuid.UUID  `gorm:"type:uuid;not null;index"`
	Name            string     `gorm:"type:text;not null;default:'Основной счет'"`
	BIC             string     `gorm:"type:text;not null;check:bic ~ '^[0-9]{9}$'"`
	BankName        string     `gorm:"type:text;not null"`
	BankAddress     string     `gorm:"type:text;not null"`
	CorrAccount     string     `gorm:"type:text;not null;check:correspondent_account ~ '^[0-9]{20}$'"`
	PaymentAccount  string     `gorm:"type:text;not null;check:payment_account ~ '^[0-9]{20}$'"`
	Currency        string     `gorm:"type:text;not null;default:'RUB'"`
	Comment         string     `gorm:"type:text;default:''"`
	IsPrimary       bool       `gorm:"type:boolean;not null;default:false"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;default:now();not null"`
	UpdatedAt       time.Time  `gorm:"type:timestamptz;default:now();not null"`
	DeletedAt       *time.Time `gorm:"type:timestamptz;default:NULL"`
}

func (BankAccount) TableName() string {
	return "bank_accounts"
}

// После создания записи или обновления, если is_primary=true,
// остальные счета этого юрлица должны стать is_primary=false.
func (ba *BankAccount) AfterSave(tx *gorm.DB) (err error) {
	if ba.IsPrimary {
		return tx.Model(&BankAccount{}).
			Where("legal_entity_uuid = ? AND uuid != ?", ba.LegalEntityUUID, ba.UUID).
			Update("is_primary", false).Error
	}
	return nil
}
