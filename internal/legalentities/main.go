package legalentities

import (
	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateLegalEntity(legalEntity *domain.LegalEntity) error {
	return s.repo.CreateLegalEntity(legalEntity)
}

func (s *Service) GetAllLegalEntities() ([]domain.LegalEntity, error) {
	return s.repo.GetAllLegalEntities()
}

func (s *Service) UpdateLegalEntity(legalEntity domain.LegalEntity) error {
	return s.repo.UpdateLegalEntity(legalEntity)
}

func (s *Service) DeleteLegalEntity(legailEntityUUID uuid.UUID) error {
	return s.repo.DeleteLegalEntity(legailEntityUUID)
}

func (s *Service) CreateBankAccount(bankAccount domain.BankAccount) error {
	return s.repo.CreateBankAccount(bankAccount)
}

func (s *Service) GetAllBankAccounts() ([]domain.BankAccount, error) {
	return s.repo.GetAllBankAccounts()
}

func (s *Service) UpdateBankAccount(bankAccount domain.BankAccount) error {
	return s.repo.UpdateBankAccount(bankAccount)
}

func (s *Service) DeleteBankAccount(bankAccountUUID uuid.UUID) error {
	return s.repo.DeleteBankAccount(bankAccountUUID)
}

func (s *Service) GetLegalEntityWithAccounts(entityUUID uuid.UUID) (*domain.LegalEntity, error) {
	return s.repo.GetLegalEntityWithAccounts(entityUUID)
}

// Добавляем в Service (legalentities/main.go).
func (s *Service) GetBankAccountByUUID(bankAccountUUID uuid.UUID) (*domain.BankAccount, error) {
	return s.repo.GetBankAccountByUUID(bankAccountUUID)
}
