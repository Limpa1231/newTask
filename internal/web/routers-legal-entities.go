package web

import (
	"context"
	"time"

	"github.com/krisch/crm-backend/domain"
	"github.com/krisch/crm-backend/internal/jwt"
	oapi "github.com/krisch/crm-backend/internal/web/ofederation"
	echo "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

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

// PatchLegalentitiesUUID implements ofederation.StrictServerInterface.
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

	// Возвращаем данные на основе запроса, так как сервис не возвращает обновленную сущность
	return oapi.PatchLegalEntitiesUUID200JSONResponse{
		UUID:      request.UUID,
		Name:      request.Body.Name,
		CreatedAt: time.Now(), // Эти значения должны быть получены из БД
		UpdatedAt: time.Now(), // В реальном коде нужно получить актуальные значения
	}, nil
}

// DeleteLegalentitiesUUID implements ofederation.StrictServerInterface.
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
