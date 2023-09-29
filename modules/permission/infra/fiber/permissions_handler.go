package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/permission/app/query"
	"github.com/oechsler-it/identity/modules/permission/domain"
)

type permissionResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionsHandler struct {
	*fiber.App
	// ---
	FindAll cqrs.QueryHandler[query.FindAll, []*domain.Permission]
}

func UsePermissionsHandler(handler *PermissionsHandler) {
	permission := handler.Group("/permission")
	permission.Get("/", handler.get)
}

// @Summary	Get all permissions
// @Produce	json
// @Success	200	{object}	permissionResponse
// @Failure	500
// @Router		/permission [get]
// @Tags		Permission
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
