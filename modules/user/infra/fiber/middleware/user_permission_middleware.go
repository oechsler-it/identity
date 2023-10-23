package middleware

import (
	"github.com/gofiber/fiber/v2"
	tokenDomain "github.com/oechsler-it/identity/modules/token/domain"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type UserPermissionMiddleware struct {
}

func (e *UserPermissionMiddleware) Has(permission domain.Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if _, ok := ctx.Locals("authorized").(struct{}); ok {
			return ctx.Next()
		}

		// If the user is authenticated by a token, we skip the user permission check,
		// because the permission set of a token can be more strict than the user permission set.
		if _, ok := ctx.Locals("token").(*tokenDomain.Token); ok {
			return ctx.Next()
		}

		user, ok := ctx.Locals("user").(*domain.User)
		if !ok {
			return ctx.Next()
		}

		if permission == domain.PermissionNone {
			ctx.Locals("authorized", struct{}{})

			return ctx.Next()
		}

		if !user.HasPermissionAkinTo(permission) {
			return ctx.Next()
		}

		ctx.Locals("authorized", struct{}{})

		return ctx.Next()
	}
}
