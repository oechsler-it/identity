package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
)

type SessionMiddleware struct {
	*fiber.App
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Session]
}

func UseSessionMiddleware(middleware *SessionMiddleware) {
	middleware.Use(middleware.handle)
}

func (e *SessionMiddleware) handle(ctx *fiber.Ctx) error {
	sessionIdCookie := ctx.Cookies("session_id")
	sessionId, err := uuid.FromString(sessionIdCookie)
	if err != nil {
		return ctx.Next()
	}

	ctx.Locals("session_id", domain.SessionId(sessionId))

	session, err := e.FindById.Handle(ctx.Context(), query.FindById{
		Id: domain.SessionId(sessionId),
	})
	if err != nil {
		return ctx.Next()
	}

	ctx.Locals("session", session)

	return ctx.Next()
}
