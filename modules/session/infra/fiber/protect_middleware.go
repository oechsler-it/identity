package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
)

type ProtectMiddleware struct {
	VerifyActive cqrs.CommandHandler[command.VerifyActive]
}

func (e *ProtectMiddleware) Handle(ctx *fiber.Ctx) error {
	sessionIdCookie := ctx.Cookies("session_id")
	if sessionIdCookie == "" {
		return fiber.ErrUnauthorized
	}

	sessionId, err := uuid.FromString(sessionIdCookie)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if err := e.VerifyActive.Handle(ctx.Context(), command.VerifyActive{
		Id: domain.SessionId(sessionId),
	}); err != nil {
		return fiber.ErrUnauthorized
	}

	return ctx.Next()
}
