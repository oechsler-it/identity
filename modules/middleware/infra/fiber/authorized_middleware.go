package fiber

import "github.com/gofiber/fiber/v2"

type AuthorizedMiddleware struct {
}

func (e *AuthorizedMiddleware) Handle(ctx *fiber.Ctx) error {
	if _, ok := ctx.Locals("authorized").(struct{}); !ok {
		return fiber.ErrForbidden
	}
	return ctx.Next()
}
