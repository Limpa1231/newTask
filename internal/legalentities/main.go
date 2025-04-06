package legalentities

import (
	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
)

type LegalEntitiesRepository interface {
	CreateLegalEntity(legalEntity *domain.LegalEntity) error
	GetAllLegalEntities() ([]domain.LegalEntity, error)
	UpdateLegalEntity(legalEntity domain.LegalEntity) error
	DeleteLegalEntity(legalEntityUUID uuid.UUID) error
}
type Service struct {
	repo LegalEntitiesRepository
}

func NewService(repo LegalEntitiesRepository) *Service {
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
