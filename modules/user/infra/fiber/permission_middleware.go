package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type PermissionMiddleware struct {
	VerifyHasPermission cqrs.CommandHandler[command.VerifyHasPermission]
}

func (e *PermissionMiddleware) Has(permission domain.Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*domain.User)
		if !ok {
			return fiber.ErrForbidden
		}

		if err := e.VerifyHasPermission.Handle(ctx.Context(), command.VerifyHasPermission{
			Id:         user.Id,
			Permission: permission,
		}); err != nil {
			return fiber.ErrForbidden
		}

		return ctx.Next()
	}
}
