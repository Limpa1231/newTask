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
