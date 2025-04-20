package web

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
	"github.com/krisch/crm-backend/internal/jwt"
	oapi "github.com/krisch/crm-backend/internal/web/ofederation"
	echo "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type LegalEntitiesService interface {
	CreateLegalEntity(legalEntity *domain.LegalEntity) error
	GetAllLegalEntities() ([]domain.LegalEntity, error)
	UpdateLegalEntity(legalEntity domain.LegalEntity) error
	DeleteLegalEntity(legalEntityUUID uuid.UUID) error
}

func initOpenAPILegalEntityRouters(a *Web, e *echo.Echo) {
	logrus.WithField("route", "oFederation").Debug("routes initialization")

	midlewares := []oapi.StrictMiddlewareFunc{
		ValidateStructMiddeware,
		AuthMiddeware(a.app, []string{}),
	}

	handlers := oapi.NewStrictHandler(a, midlewares)
	oapi.RegisterHandlers(e, handlers)
}

func (a *Web) PostLegalEntities(ctx context.Context, request oapi.PostLegalEntitiesRequestObject) (oapi.PostLegalEntitiesResponseObject, error) {
	_, ok := ctx.Value(claimsKey).(jwt.Claims)
	if !ok {
		return nil, ErrInvalidAuthHeader
	}

	// Создаем временный объект с полями
	entity := domain.LegalEntity{
		Name: request.Body.Name,
		// UUID и даты будут установлены в репозитории
	}

	// Передаем указатель, чтобы получить изменения из репозитория
	if err := a.app.LegalEntitiesService.CreateLegalEntity(&entity); err != nil {
		return nil, err
	}

	// Теперь entity содержит все обновленные поля
	return oapi.PostLegalEntities201JSONResponse{
		UUID:      entity.UUID,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

func (a *Web) GetLegalEntities(ctx context.Context, _ oapi.GetLegalEntitiesRequestObject) (oapi.GetLegalEntitiesResponseObject, error) {
	_, ok := ctx.Value(claimsKey).(jwt.Claims)
	if !ok {
		return nil, ErrInvalidAuthHeader
	}

	items, err := a.app.LegalEntitiesService.GetAllLegalEntities()
	if err != nil {
		return nil, err
	}

	legalEntities := make([]oapi.LegalEntitiesDTO, len(items))
	for i, entity := range items {
		legalEntities[i] = oapi.LegalEntitiesDTO{
			UUID:      entity.UUID,
			Name:      entity.Name,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		}
	}

	return oapi.GetLegalEntities200JSONResponse(legalEntities), nil
}

func (a *Web) PatchLegalEntitiesUUID(ctx context.Context, request oapi.PatchLegalEntitiesUUIDRequestObject) (oapi.PatchLegalEntitiesUUIDResponseObject, error) {
	if _, ok := ctx.Value(claimsKey).(jwt.Claims); !ok {
		return nil, ErrInvalidAuthHeader
	}

	// Предполагаем, что сервис возвращает обновленную сущность
	if err := a.app.LegalEntitiesService.UpdateLegalEntity(domain.LegalEntity{
		UUID: request.UUID,
		Name: request.Body.Name,
	}); err != nil {
		return nil, err
	}

	return oapi.PatchLegalEntitiesUUID200JSONResponse{
		UUID:      request.UUID,
		Name:      request.Body.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (a *Web) DeleteLegalEntitiesUUID(ctx context.Context, request oapi.DeleteLegalEntitiesUUIDRequestObject) (oapi.DeleteLegalEntitiesUUIDResponseObject, error) {
	_, ok := ctx.Value(claimsKey).(jwt.Claims)
	if !ok {
		return nil, ErrInvalidAuthHeader
	}

	err := a.app.LegalEntitiesService.DeleteLegalEntity(request.UUID)
	if err != nil {
		return nil, err
	}

	return oapi.DeleteLegalEntitiesUUID204Response{}, nil
}

func (a *Web) GetBankAccounts(ctx context.Context, _ oapi.GetBankAccountsRequestObject) (oapi.GetBankAccountsResponseObject, error) {
	_, ok := ctx.Value(claimsKey).(jwt.Claims)
	if !ok {
		return nil, ErrInvalidAuthHeader
	}

	accounts, err := a.app.LegalEntitiesService.GetAllBankAccounts()
	if err != nil {
		return nil, err
	}

	accountsDTO := make([]oapi.BankAccountDTO, len(accounts))
	for i := range accounts {
		account := accounts[i]
		var legalEntityUUID *uuid.UUID
		if account.LegalEntityUUID != uuid.Nil {
			legalEntityUUID = &account.LegalEntityUUID
		}

		accountsDTO[i] = oapi.BankAccountDTO{
			Uuid:            account.UUID,
			LegalEntityUuid: legalEntityUUID,
			Bic:             account.BIC,
			Bank:            &account.BankName,
			Address:         &account.BankAddress,
			CorrAccount:     &account.CorrAccount,
			CurrentAccount:  account.PaymentAccount,
			Currency:        account.Currency,
			Comment:         &account.Comment,
			CreatedAt:       &account.CreatedAt,
			UpdatedAt:       &account.UpdatedAt,
		}
	}
	return oapi.GetBankAccounts200JSONResponse(accountsDTO), nil
}

func (a *Web) PostBankAccounts(ctx context.Context, request oapi.PostBankAccountsRequestObject) (oapi.PostBankAccountsResponseObject, error) {
	// Проверка аутентификации
	_, ok := ctx.Value(claimsKey).(jwt.Claims)
	if !ok {
		return nil, ErrInvalidAuthHeader
	}

	// Проверка обязательных полей
	if request.Body.LegalEntityUuid == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "legal_entity_uuid is required")
	}
	if request.Body.Bic == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "bic is required")
	}
	if request.Body.CurrentAccount == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "current_account is required")
	}
	if request.Body.Currency == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "currency is required")
	}

	// Обработка опциональных полей
	var bankName, address, corrAccount string
	comment := "Нет Комментария"
	if request.Body.Bank != nil {
		bankName = *request.Body.Bank
	}
	if request.Body.Address != nil {
		address = *request.Body.Address
	}
	if request.Body.CorrAccount != nil {
		corrAccount = *request.Body.CorrAccount
	}
	if request.Body.Comment != nil {
		comment = *request.Body.Comment
	}

	account := domain.BankAccount{
		UUID:            uuid.New(),
		LegalEntityUUID: *request.Body.LegalEntityUuid,
		BIC:             request.Body.Bic,
		BankName:        bankName,
		BankAddress:     address,
		CorrAccount:     corrAccount,
		PaymentAccount:  request.Body.CurrentAccount,
		Currency:        request.Body.Currency,
		Comment:         comment, // Дефолтное значение
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := a.app.LegalEntitiesService.CreateBankAccount(account); err != nil {
		return nil, err
	}

	// Подготовка ответа
	response := oapi.PostBankAccounts201JSONResponse{
		Uuid:            account.UUID,
		LegalEntityUuid: &account.LegalEntityUUID,
		Bic:             account.BIC,
		CurrentAccount:  account.PaymentAccount,
		Currency:        account.Currency,
	}

	// Добавление опциональных полей в ответ
	if account.BankName != "" {
		response.Bank = &account.BankName
	}
	if account.BankAddress != "" {
		response.Address = &account.BankAddress
	}
	if account.CorrAccount != "" {
		response.CorrAccount = &account.CorrAccount
	}
	if account.Comment != "" {
		response.Comment = &account.Comment
	}

	createdAt := account.CreatedAt
	updatedAt := account.UpdatedAt
	response.CreatedAt = &createdAt
	response.UpdatedAt = &updatedAt

	return response, nil
}

func (a *Web) PatchBankAccountsUUID(ctx context.Context, request oapi.PatchBankAccountsUUIDRequestObject) (oapi.PatchBankAccountsUUIDResponseObject, error) {
	// Проверка обязательных полей
	if request.Body == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "request body is required")
	}

	// Получаем текущий аккаунт
	currentAccount, err := a.app.LegalEntitiesService.GetBankAccountByUUID(request.UUID)
	if err != nil {
		return nil, err
	}

	// Подготавливаем обновления
	updates := domain.BankAccount{
		UUID:      request.UUID,
		UpdatedAt: time.Now(),
	}

	// Применяем изменения только к переданным полям
	if request.Body.Bank != nil {
		updates.BankName = *request.Body.Bank
	} else {
		updates.BankName = currentAccount.BankName
	}

	if request.Body.Bic != "" {
		updates.BIC = request.Body.Bic
	} else {
		updates.BIC = currentAccount.BIC
	}

	// Аналогично для остальных полей...
	if request.Body.Address != nil {
		updates.BankAddress = *request.Body.Address
	}
	if request.Body.CorrAccount != nil {
		updates.CorrAccount = *request.Body.CorrAccount
	}
	if request.Body.CurrentAccount != "" {
		updates.PaymentAccount = request.Body.CurrentAccount
	}
	if request.Body.Currency != "" {
		updates.Currency = request.Body.Currency
	}
	if request.Body.Comment != nil {
		updates.Comment = *request.Body.Comment
	}

	if err := a.app.LegalEntitiesService.UpdateBankAccount(updates); err != nil {
		return nil, err
	}

	// Возвращаем обновлённый аккаунт
	updatedAccount, err := a.app.LegalEntitiesService.GetBankAccountByUUID(request.UUID)
	if err != nil {
		return nil, err
	}

	return oapi.PatchBankAccountsUUID200JSONResponse{
		Uuid:            updatedAccount.UUID,
		LegalEntityUuid: &updatedAccount.LegalEntityUUID,
		Bic:             updatedAccount.BIC,
		Bank:            &updatedAccount.BankName,
		Address:         &updatedAccount.BankAddress,
		CorrAccount:     &updatedAccount.CorrAccount,
		CurrentAccount:  updatedAccount.PaymentAccount,
		Currency:        updatedAccount.Currency,
		Comment:         &updatedAccount.Comment,
		CreatedAt:       &updatedAccount.CreatedAt,
		UpdatedAt:       &updatedAccount.UpdatedAt,
	}, nil
}

func (a *Web) DeleteBankAccountsUUID(ctx context.Context, request oapi.DeleteBankAccountsUUIDRequestObject) (oapi.DeleteBankAccountsUUIDResponseObject, error) {
	_, ok := ctx.Value(claimsKey).(jwt.Claims)
	if !ok {
		return nil, ErrInvalidAuthHeader
	}

	if err := a.app.LegalEntitiesService.DeleteBankAccount(request.UUID); err != nil {
		return nil, err
	}

	return oapi.DeleteBankAccountsUUID204Response{}, nil
}
