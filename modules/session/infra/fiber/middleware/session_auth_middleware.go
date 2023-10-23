package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
)

type SessionAuthMiddleware struct {
	VerifyActive cqrs.CommandHandler[command.VerifyActive]
}

func (e *SessionAuthMiddleware) Handle(ctx *fiber.Ctx) error {
	if _, ok := ctx.Locals("authenticated").(struct{}); ok {
		return ctx.Next()
	}

	sessionId, ok := ctx.Locals("session_id").(domain.SessionId)
	if !ok {
		return ctx.Next()
	}

	deviceId, ok := ctx.Locals("device_id").(domain.DeviceId)
	if !ok {
		return ctx.Next()
	}

	if err := e.VerifyActive.Handle(ctx.Context(), command.VerifyActive{
		Id:       sessionId,
		DeviceId: deviceId,
	}); err != nil {
		return ctx.Next()
	}

	ctx.Locals("authenticated", struct{}{})

	return ctx.Next()
}
