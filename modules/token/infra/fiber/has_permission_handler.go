package fiber

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
	"github.com/sirupsen/logrus"
)

type HasPermissionHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	Has cqrs.CommandHandler[command.VerifyHasPermission]
}

func UseHasPermissionHandler(handler *HasPermissionHandler) {
	token := handler.Group("/token/:id")
	has := token.Group("/has")
	has.Get("/:permission", handler.get)
}

// @Summary	Verify if a token has a permission
// @Accept		text/plain
// @Produce	text/plain
// @Param		id			path	string	true	"Id of the token"
// @Param		permission	path	string	true	"Name of the permission"
// @Success	204
// @Failure	400
// @Failure	403
// @Failure	404
// @Failure	500
// @Router		/token/{id}/has/{permission} [get]
// @Tags		Token
func (e *HasPermissionHandler) get(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	idPartial, err := domain.NewTokenIdPartialFromString(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid token id")
	}

	permission := ctx.Params("permission")

	if err := e.Has.Handle(ctx.Context(), command.VerifyHasPermission{
		Id:         idPartial,
		Permission: domain.Permission(permission),
	}); err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrTokenIsExpired) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrTokenDoesNotHavePermission) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
