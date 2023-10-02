package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	sessionDomain "github.com/oechsler-it/identity/modules/session/domain"
	tokenDomain "github.com/oechsler-it/identity/modules/token/domain"
	"github.com/oechsler-it/identity/modules/user/app/query"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
)

type UserMiddleware struct {
	FindById cqrs.QueryHandler[query.FindByIdentifier, *domain.User]
}

func (e *UserMiddleware) Handle(ctx *fiber.Ctx) error {
	userId := e.findBySession(ctx)
	if userId == domain.UserId(uuid.Nil) {
		userId = e.findByToken(ctx)
	}
	if userId == domain.UserId(uuid.Nil) {
		return ctx.Next()
	}

	ctx.Locals("user_id", userId)

	user, err := e.FindById.Handle(ctx.Context(), query.FindByIdentifier{
		Identifier: uuid.UUID(userId).String(),
	})
	if err != nil {
		return ctx.Next()
	}

	ctx.Locals("user", user)

	return ctx.Next()
}

func (e *UserMiddleware) findBySession(ctx *fiber.Ctx) domain.UserId {
	session, ok := ctx.Locals("session").(*sessionDomain.Session)
	if !ok {
		return domain.UserId(uuid.Nil)
	}

	return domain.UserId(session.OwnedBy.UserId)
}

func (e *UserMiddleware) findByToken(ctx *fiber.Ctx) domain.UserId {
	token, ok := ctx.Locals("token").(*tokenDomain.Token)
	if !ok {
		return domain.UserId(uuid.Nil)
	}

	return domain.UserId(token.OwnedBy.UserId)
}
