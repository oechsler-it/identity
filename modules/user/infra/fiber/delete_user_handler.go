package fiber

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type DeleteUserHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	RenewMiddleware          *sessionFiber.RenewMiddleware
	ProtectSessionMiddleware *sessionFiber.ProtectSessionMiddleware
	UserMiddleware           *UserMiddleware
	UserPermissionMiddleware *UserPermissionMiddleware
	// ---
	Delete cqrs.CommandHandler[command.Delete]
}

func UseDeleteUserHandler(handler *DeleteUserHandler) {
	user := handler.Group("/user")
	user.Delete("/:id",
		handler.RenewMiddleware.Handle,
		handler.ProtectSessionMiddleware.Handle,
		handler.UserMiddleware.Handle,
		handler.UserPermissionMiddleware.Has("all:user:delete"),
		handler.delete)
}

// @Summary	Delete a user
// @Produce	text/plain
// @Param		id	path	string	true	"Id of the user"
// @Success	204
// @Failure	401
// @Failure	403
// @Failure	404
// @Failure	500
// @Router		/user/{id} [delete]
// @Tags		User
func (e *DeleteUserHandler) delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := uuid.FromString(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := e.Delete.Handle(ctx.Context(), command.Delete{
		Id: domain.UserId(id),
	}); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrCanNotDeleteLastUser) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithField("id", id.String()).
		Info("User deleted")

	return ctx.SendStatus(fiber.StatusNoContent)
}
