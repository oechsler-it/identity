package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/sirupsen/logrus"
)

type DeleteHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	ProtectMiddleware *sessionFiber.ProtectMiddleware
	// ---
	Delete cqrs.CommandHandler[command.Delete]
}

func UseDeleteHandler(handler *DeleteHandler) {
	delete := handler.Group("/permission")
	delete.Use(handler.ProtectMiddleware.Handle)
	delete.Delete("/:name", handler.delete)
}

//	@Summary	Delete a permission
//	@Accept		text/plain
//	@Produce	text/plain
//	@Param		name	path	string	true	"Name of the permission"
//	@Success	204
//	@Failure	401
//	@Failure	404
//	@Router		/permission/{name} [delete]
//	@Tags		Permission
func (e *DeleteHandler) delete(ctx *fiber.Ctx) error {
	name := ctx.Params("name")

	if err := e.Delete.Handle(ctx.Context(), command.Delete{
		Name: name,
	}); err != nil {
		if errors.Is(err, domain.ErrPermissionNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithFields(logrus.Fields{
		"name": name,
	}).Info("Permission deleted")

	return ctx.SendStatus(fiber.StatusNoContent)
}
