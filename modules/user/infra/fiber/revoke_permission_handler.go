package fiber

import (
	"errors"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	permissionDomain "github.com/oechsler-it/identity/modules/permission/domain"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type RevokePermissionHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	RenewMiddleware      *sessionFiber.RenewMiddleware
	ProtectMiddleware    *sessionFiber.ProtectSessionMiddleware
	UserMiddleware       *UserMiddleware
	PermissionMiddleware *UserPermissionMiddleware
	// ---
	Revoke cqrs.CommandHandler[command.RevokePermission]
}

func UseRevokePermissionHandler(handler *RevokePermissionHandler) {
	user := handler.Group("/user/:id")
	revoke := user.Group("/revoke")
	revoke.Use(handler.RenewMiddleware.Handle)
	revoke.Use(handler.ProtectMiddleware.Handle)
	revoke.Use(handler.UserMiddleware.Handle)
	revoke.Use(handler.PermissionMiddleware.Has("all:user:permission:revoke"))
	revoke.Delete("/:permission", handler.delete)
}

// @Summary	Revoke a permission from a user
// @Accept		text/plain
// @Produce	text/plain
// @Param		id			path	string	true	"Id of the user"
// @Param		permission	path	string	true	"Name of the permission"
// @Success	204
// @Failure	400
// @Failure	401
// @Failure	404
// @Failure	500
// @Router		/user/{id}/revoke/{permission} [delete]
// @Tags		User
func (e *RevokePermissionHandler) delete(ctx *fiber.Ctx) error {
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

	if err := e.Revoke.Handle(ctx.Context(), command.RevokePermission{
		Id:         domain.UserId(id),
		Permission: domain.Permission(permissionUnescaped),
	}); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, permissionDomain.ErrPermissionNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrUserDoesNotHavePermission) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithFields(logrus.Fields{
		"id":         id,
		"permission": permission,
	}).Info("Permission revoked")

	return ctx.SendStatus(fiber.StatusNoContent)
}
