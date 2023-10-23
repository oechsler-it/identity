package fiber

import "github.com/gofiber/fiber/v2"

type AuthenticatedMiddleware struct {
}

func (e *AuthenticatedMiddleware) Handle(ctx *fiber.Ctx) error {
	if _, ok := ctx.Locals("authenticated").(struct{}); !ok {
		return fiber.ErrUnauthorized
	}
	return ctx.Next()
}
