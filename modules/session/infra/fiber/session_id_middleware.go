package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
)

type SessionIdMiddleware struct {
	*fiber.App
}

func UseSessionIdMiddleware(middleware *SessionIdMiddleware) {
	middleware.Use(middleware.handle)
}

func (e *SessionIdMiddleware) handle(ctx *fiber.Ctx) error {
	sessionIdCookie := ctx.Cookies("session_id")
	sessionId, err := uuid.FromString(sessionIdCookie)
	if err != nil {
		return ctx.Next()
	}

	ctx.Locals("session_id", domain.SessionId(sessionId))

	return ctx.Next()
}
