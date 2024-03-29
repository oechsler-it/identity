package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
	userFiberMiddleware "github.com/oechsler-it/identity/modules/user/infra/fiber/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type DeleteMeHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	UserMiddleware *userFiberMiddleware.UserMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	// ---
	Delete cqrs.CommandHandler[command.Delete]
}

func UseDeleteMeHandler(handler *DeleteMeHandler) {
	user := handler.Group("/user")
	user.Delete("/me",
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.UserMiddleware.Handle,
		// ---
		handler.AuthenticatedMiddleware.Handle,
		// ---
		handler.delete)
}

// @Summary	Delete the current user
// @Produce	text/plain
// @Success	204
// @Failure	401
// @Failure	403
// @Failure	500
// @Router		/user/me [delete]
// @Tags		User
func (e *DeleteMeHandler) delete(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(*domain.User)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := e.Delete.Handle(ctx.Context(), command.Delete{
		Id: user.Id,
	}); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrCanNotDeleteLastUser) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithField("id", uuid.UUID(user.Id).String()).
		Info("User (self) deleted")

	return ctx.SendStatus(fiber.StatusNoContent)
}
