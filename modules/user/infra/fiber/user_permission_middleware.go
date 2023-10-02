package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type UserPermissionMiddleware struct {
}

func (e *UserPermissionMiddleware) Has(permission domain.Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if _, ok := ctx.Locals("authorized").(struct{}); ok {
			return ctx.Next()
		}

		user, ok := ctx.Locals("user").(*domain.User)
		if !ok {
			return fiber.ErrForbidden
		}

		if !user.HasPermissionAkinTo(permission) {
			return fiber.ErrForbidden
		}

		ctx.Locals("authorized", struct{}{})

		return ctx.Next()
	}
}
