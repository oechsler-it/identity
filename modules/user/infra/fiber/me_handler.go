package fiber

import (
	"github.com/gofiber/fiber/v2"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
	"github.com/oechsler-it/identity/modules/user/domain"
	userFiberMiddleware "github.com/oechsler-it/identity/modules/user/infra/fiber/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type MeHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	TokenAuthMiddleware *tokenFiberMiddleware.TokenAuthMiddleware
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	UserMiddleware *userFiberMiddleware.UserMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
}

func UseMeHandler(handler *MeHandler) {
	user := handler.Group("/user")
	user.Get("/me",
		handler.TokenAuthMiddleware.Handle,
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.UserMiddleware.Handle,
		// ---
		handler.AuthenticatedMiddleware.Handle,
		// ---
		handler.get)
}

// @Summary	Get information about the current user
// @Produce	json
// @Success	200	{object}	userResponse
// @Failure	401
// @Failure	500
// @Router		/user/me [get]
// @Security	TokenAuth
// @Tags		User
func (e *MeHandler) get(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(*domain.User)
	if !ok {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(userResponse{
		Id:           uuid.UUID(user.Id).String(),
		RegisteredAt: user.CreatedAt.String(),
	})
}
