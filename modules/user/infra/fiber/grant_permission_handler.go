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
	userFiberMiddleware "github.com/oechsler-it/identity/modules/user/infra/fiber/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	permissionDomain "github.com/oechsler-it/identity/modules/permission/domain"
)

type GrantPermissionHandler struct {
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
	Grant cqrs.CommandHandler[command.GrantPermission]
}

func UseGrantPermissionHandler(handler *GrantPermissionHandler) {
	user := handler.Group("/user/:id")
	grant := user.Group("/grant")
	grant.Post("/:permission",
		handler.TokenAuthMiddleware.Handle,
		handler.TokenPermissionMiddleware.Has("all:user:permission:grant"),
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.UserMiddleware.Handle,
		handler.UserPermissionMiddleware.Has("all:user:permission:grant"),
		// ---
		handler.AuthenticatedMiddleware.Handle,
		handler.AuthorizedMiddleware.Handle,
		// ---
		handler.post)
}

// @Summary	Grant a permission to a user
// @Accept		text/plain
// @Produce	text/plain
// @Param		id			path	string	true	"Id of the user"
// @Param		permission	path	string	true	"Name of the permission"
// @Success	204
// @Failure	400
// @Failure	401
// @Failure	403
// @Failure	404
// @Failure	500
// @Router		/user/{id}/grant/{permission} [post]
// @Security	TokenAuth
// @Tags		User
func (e *GrantPermissionHandler) post(ctx *fiber.Ctx) error {
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

	if err := e.Grant.Handle(ctx.Context(), command.GrantPermission{
		Id:         domain.UserId(id),
		Permission: domain.Permission(permissionUnescaped),
	}); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, permissionDomain.ErrPermissionNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrUserAlreadyHasPermission) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithFields(logrus.Fields{
		"id":         id,
		"permission": permission,
	}).Info("Permission granted")

	return ctx.SendStatus(fiber.StatusNoContent)
}
