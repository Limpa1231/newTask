package legalentities

import (
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
	var dbEntities []struct {
		UUID      uuid.UUID
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	// 1. Загружаем юрлица
	err := r.gorm.DB.Table("legal_entities").
		Where("deleted_at IS NULL").
		Find(&dbEntities).Error
	if err != nil {
		return nil, err
	}

	// 2. Загружаем все банковские счета для этих юрлиц
	var accounts []domain.BankAccount
	err = r.gorm.DB.
		Where("legal_entity_uuid IN (SELECT uuid FROM legal_entities WHERE deleted_at IS NULL)").
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	// 3. Группируем счета по legal_entity_uuid
	accountsByLegalEntity := make(map[uuid.UUID][]domain.BankAccount)
	for _, acc := range accounts {
		accountsByLegalEntity[acc.LegalEntityUUID] = append(accountsByLegalEntity[acc.LegalEntityUUID], acc)
	}

	// 4. Собираем результат
	entities := make([]domain.LegalEntity, len(dbEntities))
	for i, e := range dbEntities {
		entities[i] = domain.LegalEntity{
			UUID:         e.UUID,
			Name:         e.Name,
			CreatedAt:    e.CreatedAt,
			UpdatedAt:    e.UpdatedAt,
			BankAccounts: accountsByLegalEntity[e.UUID], // Привязываем счета
		}
	}

	return entities, nil
}

func (r *Repository) CreateLegalEntity(legalEntity *domain.LegalEntity) error {
	// Используем временную структуру для создания
	type LegalEntityDB struct {
		UUID      uuid.UUID
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	dbEntity := LegalEntityDB{
		UUID:      uuid.New(),
		Name:      legalEntity.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Создаем запись в БД
	err := r.gorm.DB.Table("legal_entities").Create(&dbEntity).Error
	if err != nil {
		return err
	}

	// Заполняем исходную структуру
	legalEntity.UUID = dbEntity.UUID
	legalEntity.CreatedAt = dbEntity.CreatedAt
	legalEntity.UpdatedAt = dbEntity.UpdatedAt

	return nil
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

func (r *Repository) GetAllBankAccounts() ([]domain.BankAccount, error) {
	var accounts []domain.BankAccount
	result := r.gorm.DB.Where("deleted_at IS NULL").Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}

func (r *Repository) CreateBankAccount(bankAccount domain.BankAccount) error {
	// Аналогичная валидация и настройка как в CreateLegalEntity
	if bankAccount.UUID == uuid.Nil {
		bankAccount.UUID = uuid.New()
	}

	now := time.Now()
	bankAccount.CreatedAt = now
	bankAccount.UpdatedAt = now

	// Добавить валидацию для БИК, номеров счетов и т.д.

	return r.gorm.DB.Create(&bankAccount).Error
}

func (r *Repository) UpdateBankAccount(bankAccount domain.BankAccount) error {
	result := r.gorm.DB.Model(&domain.BankAccount{}).
		Where("uuid = ? AND deleted_at IS NULL", bankAccount.UUID).
		Updates(map[string]interface{}{
			"bic":       bankAccount.BIC,
			"bank_name": bankAccount.BankName,
			// ... другие поля
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return dto.NotFoundErr("Банковский счет не найден")
	}

	return nil
}

func (r *Repository) DeleteBankAccount(bankAccountUUID uuid.UUID) error {
	// Мягкое удаление банковского счета (установка deleted_at)
	result := r.gorm.DB.Model(&domain.BankAccount{}).
		Where("uuid = ? AND deleted_at IS NULL", bankAccountUUID).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return dto.NotFoundErr("Банковский счет не найден или уже удален")
	}

	return nil
}

func (r *Repository) GetLegalEntityWithAccounts(entityUUID uuid.UUID) (*domain.LegalEntity, error) {
	// 1. Получаем только данные юрлица
	var legalEntity struct {
		UUID      uuid.UUID
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	err := r.gorm.DB.Table("legal_entities").
		Where("uuid = ? AND deleted_at IS NULL", entityUUID).
		First(&legalEntity).Error
	if err != nil {
		return nil, err
	}

	// 2. Получаем счета отдельным запросом
	var accounts []domain.BankAccount
	if err := r.gorm.DB.
		Where("legal_entity_uuid = ? AND deleted_at IS NULL", entityUUID).
		Find(&accounts).Error; err != nil {
		return nil, err
	}

	// 3. Собираем результат
	result := &domain.LegalEntity{
		UUID:         legalEntity.UUID,
		Name:         legalEntity.Name,
		CreatedAt:    legalEntity.CreatedAt,
		UpdatedAt:    legalEntity.UpdatedAt,
		BankAccounts: accounts,
	}

	return result, nil
}

// Добавляем в Repository (legalentities/repository.go).
func (r *Repository) GetBankAccountByUUID(bankAccountUUID uuid.UUID) (*domain.BankAccount, error) {
	var account domain.BankAccount
	result := r.gorm.DB.Where("uuid = ? AND deleted_at IS NULL", bankAccountUUID).First(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}
