package legalentities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
	"github.com/krisch/crm-backend/dto"
	"github.com/krisch/crm-backend/pkg/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	gorm *postgres.GDB
}

var _ LegalEntitiesRepository = (*Repository)(nil)

func NewRepository(gorm *postgres.GDB) *Repository {
	return &Repository{gorm: gorm}
}

func (r *Repository) handleError(result *gorm.DB, notFoundMsg string) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return dto.NotFoundErr(notFoundMsg)
	}
	return nil
}

func (r *Repository) createEntity(entity *domain.LegalEntity, validateFunc func() error) error {
	if entity.UUID == uuid.Nil {
		entity.UUID = uuid.New()
	}

	if err := validateFunc(); err != nil {
		return err
	}

	now := time.Now()
	entity.CreatedAt = now
	entity.UpdatedAt = now

	return r.gorm.DB.Select("UUID", "Name", "CreatedAt", "UpdatedAt").Create(entity).Error
}

func (r *Repository) CreateLegalEntity(legalEntity *domain.LegalEntity) error {
	return r.createEntity(legalEntity, func() error {
		if legalEntity.Name == "" {
			return errors.New("name is required")
		}

		var count int64
		if err := r.gorm.DB.Model(&domain.LegalEntity{}).
			Where("name = ? AND deleted_at IS NULL", legalEntity.Name).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("юр. лицо с таким именем уже существует")
		}
		return nil
	})
}

func (r *Repository) GetAllLegalEntities() ([]domain.LegalEntity, error) {
	var entities []domain.LegalEntity
	result := r.gorm.DB.Where("deleted_at IS NULL").Find(&entities)
	return entities, result.Error
}

func (r *Repository) UpdateLegalEntity(legalEntity domain.LegalEntity) error {
	result := r.gorm.DB.Model(&domain.LegalEntity{}).
		Where("uuid = ? AND deleted_at IS NULL", legalEntity.UUID).
		Updates(map[string]interface{}{
			"name":       legalEntity.Name,
			"updated_at": time.Now(),
		})
	return r.handleError(result, "Юр. лицо не найдено")
}

func (r *Repository) DeleteLegalEntity(legalEntityUUID uuid.UUID) error {
	result := r.gorm.DB.Model(&domain.LegalEntity{}).
		Where("uuid = ? AND deleted_at IS NULL", legalEntityUUID).
		Update("deleted_at", time.Now())
	return r.handleError(result, "Юр. лицо не найдено")
}
