package fiber

import (
	"errors"

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
	ProtectMiddleware *sessionFiber.ProtectMiddleware
	// ---
	Grant cqrs.CommandHandler[command.GrantPermission]
}

func UseGrantPermissionHandler(handler *GrantPermissionHandler) {
	user := handler.Group("/user/:id")
	grant := user.Group("/grant")
	grant.Use(handler.ProtectMiddleware.Handle)
	grant.Post("/:permission", handler.post)
}

//	@Summary	Grant a permission to a user
//	@Accept		text/plain
//	@Produce	text/plain
//	@Param		id			path	string	true	"ID of the user"
//	@Param		permission	path	string	true	"Name of the permission"
//	@Success	204
//	@Failure	401
//	@Failure	403
//	@Failure	404
//	@Router		/user/{id}/grant/{permission} [post]
//	@Tags		User
func (e *GrantPermissionHandler) post(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := uuid.FromString(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	permission := ctx.Params("permission")

	if err := e.Grant.Handle(ctx.Context(), command.GrantPermission{
		Id:         domain.UserId(id),
		Permission: domain.Permission(permission),
	}); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, permissionDomain.ErrPermissionNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrUserAlreadyHasPermission) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithFields(logrus.Fields{
		"id":         id,
		"permission": permission,
	}).Info("Permission granted")

	return ctx.SendStatus(fiber.StatusNoContent)
}
