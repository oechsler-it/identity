package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/modules/token/domain"
)

type TokenPermissionMiddleware struct {
}

func (e *TokenPermissionMiddleware) Has(permission domain.Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if _, ok := ctx.Locals("authorized").(struct{}); ok {
			return ctx.Next()
		}

		token, ok := ctx.Locals("token").(*domain.Token)
		if !ok {
			return fiber.ErrForbidden
		}

		if !token.HasPermissionAkinTo(permission) {
			return fiber.ErrForbidden
		}

		ctx.Locals("authorized", struct{}{})

		return ctx.Next()
	}
}
