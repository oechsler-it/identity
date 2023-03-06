package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	"github.com/oechsler-it/identity/runtime"
	uuid "github.com/satori/go.uuid"
	"time"
)

type RenewMiddleware struct {
	Env *runtime.Env
	// ---
	Renew cqrs.CommandHandler[command.Renew]
}

func (e *RenewMiddleware) Handle(ctx *fiber.Ctx) error {
	sessionIdCookie := ctx.Cookies("session_id")
	if sessionIdCookie == "" {
		return ctx.Next()
	}

	sessionId, err := uuid.FromString(sessionIdCookie)
	if err != nil {
		return ctx.Next()
	}

	lifetimeInSeconds := e.Env.Int("SESSION_LIFETIME_IN_HOURS", 8) * 60 * 60

	if err := e.Renew.Handle(ctx.Context(), command.Renew{
		Id:                   domain.SessionId(sessionId),
		NewLifeTimeInSeconds: lifetimeInSeconds,
	}); err != nil {
		return ctx.Next()
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Value:   sessionId.String(),
		Path:    "/",
		Expires: time.Now().Add(time.Duration(lifetimeInSeconds) * time.Second),
	})

	return ctx.Next()
}
