package fiber

import (
	"errors"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type HasPermissionHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	TokenAuthMiddleware *tokenFiberMiddleware.TokenAuthMiddleware
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	// ---
	Has cqrs.CommandHandler[command.VerifyHasPermission]
}

func UseHasPermissionHandler(handler *HasPermissionHandler) {
	user := handler.Group("/user/:id")
	has := user.Group("/has")
	has.Get("/:permission",
		handler.TokenAuthMiddleware.Handle,
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.AuthenticatedMiddleware.Handle,
		// ---
		handler.get)
}

// @Summary	Verify if a user has a permission
// @Accept		text/plain
// @Produce	text/plain
// @Param		id			path	string	true	"Id of the user"
// @Param		permission	path	string	true	"Name of the permission"
// @Success	204
// @Failure	400
// @Failure	403
// @Failure	404
// @Failure	500
// @Router		/user/{id}/has/{permission} [get]
// @Security	TokenAuth
// @Tags		User
func (e *HasPermissionHandler) get(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := uuid.FromString(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	permission := ctx.Params("permission")
	permissionUnescaped, err := url.PathUnescape(permission)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := e.Has.Handle(ctx.Context(), command.VerifyHasPermission{
		Id:         domain.UserId(id),
		Permission: domain.Permission(permissionUnescaped),
	}); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrUserDoesNotHavePermission) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
