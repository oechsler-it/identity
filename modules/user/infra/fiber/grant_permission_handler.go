package fiber

import (
	"errors"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	permissionDomain "github.com/oechsler-it/identity/modules/permission/domain"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
)

type GrantPermissionHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	RenewMiddleware      *sessionFiber.RenewMiddleware
	ProtectMiddleware    *sessionFiber.ProtectMiddleware
	UserMiddleware       *UserMiddleware
	PermissionMiddleware *PermissionMiddleware
	// ---
	Grant cqrs.CommandHandler[command.GrantPermission]
}

func UseGrantPermissionHandler(handler *GrantPermissionHandler) {
	user := handler.Group("/user/:id")
	grant := user.Group("/grant")
	grant.Use(handler.RenewMiddleware.Handle)
	grant.Use(handler.ProtectMiddleware.Handle)
	grant.Use(handler.UserMiddleware.Handle)
	grant.Use(handler.PermissionMiddleware.Has("all:user:permission:grant"))
	grant.Post("/:permission", handler.post)
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
