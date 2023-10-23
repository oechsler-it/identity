package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/token/app/query"
	"github.com/oechsler-it/identity/modules/token/domain"
)

type TokenMiddleware struct {
	*fiber.App
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Token]
}

func UseTokenMiddleware(middleware *TokenMiddleware) {
	middleware.Use(middleware.handle)
}

func (e *TokenMiddleware) handle(ctx *fiber.Ctx) error {
	authorizationHeader := ctx.Get("Authorization")
	if authorizationHeader == "" {
		return ctx.Next()
	}

	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return ctx.Next()
	}

	tokenId := domain.TokenId(strings.TrimPrefix(authorizationHeader, "Bearer "))

	ctx.Locals("token_id", tokenId)

	token, err := e.FindById.Handle(ctx.Context(), query.FindById{
		Id: tokenId,
	})
	if err != nil {
		return ctx.Next()
	}

	ctx.Locals("token", token)

	return ctx.Next()
}
