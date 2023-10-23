package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
)

type TokenAuthMiddleware struct {
	VerifyActive cqrs.CommandHandler[command.VerifyActive]
}

func (e *TokenAuthMiddleware) Handle(ctx *fiber.Ctx) error {
	if _, ok := ctx.Locals("authenticated").(struct{}); ok {
		return ctx.Next()
	}

	tokenId, ok := ctx.Locals("token_id").(domain.TokenId)
	if !ok {
		return ctx.Next()
	}

	if err := e.VerifyActive.Handle(ctx.Context(), command.VerifyActive{
		Id: tokenId,
	}); err != nil {
		return ctx.Next()
	}

	ctx.Locals("authenticated", struct{}{})

	return ctx.Next()
}
