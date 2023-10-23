package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	"github.com/oechsler-it/identity/runtime"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type RenewMiddleware struct {
	Logger *logrus.Logger
	Env    *runtime.Env
	// ---
	Renew cqrs.CommandHandler[command.Renew]
}

func (e *RenewMiddleware) Handle(ctx *fiber.Ctx) error {
	sessionId, ok := ctx.Locals("session_id").(domain.SessionId)
	if !ok {
		return ctx.Next()
	}

	lifetimeInSeconds := e.Env.Int("SESSION_LIFETIME_IN_HOURS", 8) * 60 * 60

	if err := e.Renew.Handle(ctx.Context(), command.Renew{
		Id:                   sessionId,
		NewLifeTimeInSeconds: lifetimeInSeconds,
	}); err != nil {
		return ctx.Next()
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Value:   uuid.UUID(sessionId).String(),
		Path:    "/",
		Expires: time.Now().Add(time.Duration(lifetimeInSeconds) * time.Second),
	})

	e.Logger.WithFields(logrus.Fields{
		"session_id": uuid.UUID(sessionId).String(),
	}).Info("Session renewed")

	return ctx.Next()
}
