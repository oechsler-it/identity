package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
	userFiberMiddleware "github.com/oechsler-it/identity/modules/user/infra/fiber/middleware"
	"github.com/sirupsen/logrus"
)

type DeleteHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	TokenAuthMiddleware       *tokenFiberMiddleware.TokenAuthMiddleware
	TokenPermissionMiddleware *tokenFiberMiddleware.TokenPermissionMiddleware
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	UserMiddleware           *userFiberMiddleware.UserMiddleware
	UserPermissionMiddleware *userFiberMiddleware.UserPermissionMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	AuthorizedMiddleware    *middlewareFiber.AuthorizedMiddleware
	// ---
	Delete cqrs.CommandHandler[command.Delete]
}

func UseDeleteHandler(handler *DeleteHandler) {
	del := handler.Group("/permission")
	del.Delete("/:name",
		handler.TokenAuthMiddleware.Handle,
		handler.TokenPermissionMiddleware.Has("all:permission:delete"),
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.UserMiddleware.Handle,
		handler.UserPermissionMiddleware.Has("all:permission:delete"),
		// ---
		handler.AuthenticatedMiddleware.Handle,
		handler.AuthorizedMiddleware.Handle,
		// ---
		handler.delete)
}

//	@Summary	Delete a permission
//	@Accept		text/plain
//	@Produce	text/plain
//	@Param		name	path	string	true	"Name of the permission"
//	@Success	204
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/permission/{name} [delete]
//	@Security	TokenAuth
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
