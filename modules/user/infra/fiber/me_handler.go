package fiber

import (
	"github.com/gofiber/fiber/v2"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type MeHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	RenewMiddleware          *sessionFiber.RenewMiddleware
	ProtectSessionMiddleware *sessionFiber.ProtectSessionMiddleware
	UserMiddleware           *UserMiddleware
}

func UseMeHandler(handler *MeHandler) {
	user := handler.Group("/user")
	user.Get("/me",
		handler.RenewMiddleware.Handle,
		handler.ProtectSessionMiddleware.Handle,
		handler.UserMiddleware.Handle,
		handler.get)
}

// @Summary	Get information about the current user
// @Produce	json
// @Success	200	{object}	userResponse
// @Failure	401
// @Failure	500
// @Router		/user/me [get]
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
