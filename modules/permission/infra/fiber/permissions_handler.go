package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	"github.com/oechsler-it/identity/modules/permission/app/query"
	"github.com/oechsler-it/identity/modules/permission/domain"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
)

type permissionResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionsHandler struct {
	*fiber.App
	// ---
	TokenAuthMiddleware *tokenFiberMiddleware.TokenAuthMiddleware
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	// ---
	FindAll cqrs.QueryHandler[query.FindAll, []*domain.Permission]
}

func UsePermissionsHandler(handler *PermissionsHandler) {
	permission := handler.Group("/permission")
	permission.Get("/",
		handler.TokenAuthMiddleware.Handle,
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.AuthenticatedMiddleware.Handle,
		// ---
		handler.get)
}

//	@Summary	Get all permissions
//	@Produce	json
//	@Success	200	{object}	permissionResponse
//	@Failure	500
//	@Router		/permission [get]
//	@Security	TokenAuth
//	@Tags		Permission
func (e *PermissionsHandler) get(ctx *fiber.Ctx) error {
	permissions, err := e.FindAll.Handle(ctx.Context(), query.FindAll{})
	if err != nil {
		return err
	}

	// ---

	response := make([]permissionResponse, len(permissions))
	for i, permission := range permissions {
		response[i] = permissionResponse{
			Name:        string(permission.Name),
			Description: permission.Description,
		}
	}
	return ctx.JSON(response)
}
