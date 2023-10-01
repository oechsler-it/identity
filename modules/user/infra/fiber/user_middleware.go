package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	sessionQuery "github.com/oechsler-it/identity/modules/session/app/query"
	sessionDomain "github.com/oechsler-it/identity/modules/session/domain"
	"github.com/oechsler-it/identity/modules/user/app/query"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
)

type UserMiddleware struct {
	FindSessionById cqrs.QueryHandler[sessionQuery.FindById, *sessionDomain.Session]
	FindById        cqrs.QueryHandler[query.FindByIdentifier, *domain.User]
}

func (e *UserMiddleware) Handle(ctx *fiber.Ctx) error {
	sessionId, ok := ctx.Locals("session_id").(sessionDomain.SessionId)
	if !ok {
		return ctx.Next()
	}

	session, err := e.FindSessionById.Handle(ctx.Context(), sessionQuery.FindById{
		Id: sessionId,
	})
	if err != nil {
		return ctx.Next()
	}

	user, err := e.FindById.Handle(ctx.Context(), query.FindByIdentifier{
		Identifier: uuid.UUID(session.OwnedBy.UserId).String(),
	})
	if err != nil {
		return ctx.Next()
	}

	ctx.Locals("user_id", user.Id)
	ctx.Locals("user", user)

	return ctx.Next()
}
