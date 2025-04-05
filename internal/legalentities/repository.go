package legalentities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
	"github.com/krisch/crm-backend/dto"
	"github.com/krisch/crm-backend/pkg/postgres"
)

type Repository struct {
	gorm *postgres.GDB
}

func NewRepository(gorm *postgres.GDB) *Repository {
	return &Repository{gorm: gorm}
}

func (r *Repository) GetAllLegalEntities() ([]domain.LegalEntity, error) {
	var entities []domain.LegalEntity
	result := r.gorm.DB.Where("deleted_at IS NULL").Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

func (r *Repository) CreateLegalEntity(legalEntity *domain.LegalEntity) error {
	if legalEntity.UUID == uuid.Nil {
		legalEntity.UUID = uuid.New()
	}

	// Убедимся, что имя не пустое
	if legalEntity.Name == "" {
		return errors.New("name is required")
	}

	// Проверка уникальности имени (ваш текущий код)
	var count int64
	if err := r.gorm.DB.Model(&domain.LegalEntity{}).
		Where("name = ? AND deleted_at IS NULL", legalEntity.Name).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("юр. лицо с таким именем уже существует")
	}

	// Установка временных меток
	now := time.Now()
	legalEntity.CreatedAt = now
	legalEntity.UpdatedAt = now

	// Явно указываем какие поля сохранять
	return r.gorm.DB.Select("UUID", "Name", "CreatedAt", "UpdatedAt").Create(legalEntity).Error
}

func (r *Repository) UpdateLegalEntity(legalEntity domain.LegalEntity) error {
	// Обновляем только имя и updated_at
	result := r.gorm.DB.Model(&domain.LegalEntity{}).
		Where("uuid = ? AND deleted_at IS NULL", legalEntity.UUID).
		Updates(map[string]interface{}{
			"name":       legalEntity.Name,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return dto.NotFoundErr("Юр. лицо не найдено")
	}

	return nil
}

func (r *Repository) DeleteLegalEntity(legalEntityUUID uuid.UUID) error {
	// Мягкое удаление (установка deleted_at)
	result := r.gorm.DB.Model(&domain.LegalEntity{}).
		Where("uuid = ? AND deleted_at IS NULL", legalEntityUUID).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return dto.NotFoundErr("Юр. лицо не найдено")
	}

	return nil
}
