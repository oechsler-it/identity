package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
)

type ProtectSessionMiddleware struct {
	VerifyActive cqrs.CommandHandler[command.VerifyActive]
}

func (e *ProtectSessionMiddleware) Handle(ctx *fiber.Ctx) error {
	if _, ok := ctx.Locals("authenticated").(struct{}); ok {
		return ctx.Next()
	}

	sessionId, ok := ctx.Locals("session_id").(domain.SessionId)
	if !ok {
		return fiber.ErrUnauthorized
	}

	deviceId, ok := ctx.Locals("device_id").(domain.DeviceId)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := e.VerifyActive.Handle(ctx.Context(), command.VerifyActive{
		Id:       sessionId,
		DeviceId: deviceId,
	}); err != nil {
		return fiber.ErrUnauthorized
	}

	ctx.Locals("authenticated", struct{}{})

	return ctx.Next()
}
