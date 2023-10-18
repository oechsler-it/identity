package fiber

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/token/app/query"
	"github.com/oechsler-it/identity/modules/token/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type TokenByIdHandler struct {
	*fiber.App
	// ---
	RenewMiddleware          *sessionFiber.RenewMiddleware
	ProtectSessionMiddleware *sessionFiber.ProtectSessionMiddleware
	// ---
	FindByIdPartial cqrs.QueryHandler[query.FindByIdPartial, *domain.Token]
}

func UseTokenByIdHandler(handler *TokenByIdHandler) {
	token := handler.Group("/token")
	token.Get("/:id",
		handler.RenewMiddleware.Handle,
		handler.ProtectSessionMiddleware.Handle,
		handler.get)
}

// @Summary	Get details of a token
// @Produce	json
// @Param		id	path		string	true	"Id of the token"
// @Success	200	{object}	tokenResponse
// @Failure	401
// @Failure	404
// @Failure	500
// @Router		/token/{id} [get]
// @Tags		Token
func (e *TokenByIdHandler) get(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	idPartial, err := domain.NewTokenIdPartialFromString(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid token id")
	}

	token, err := e.FindByIdPartial.Handle(ctx.Context(), query.FindByIdPartial{
		IdPartial: idPartial,
	})
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return err
	}

	return ctx.JSON(tokenResponse{
		Id:          token.Id.GetPartial().String(),
		Description: token.Description,
		OwnedBy: tokenOwner{
			UserId: uuid.UUID(token.OwnedBy.UserId).String(),
		},
		IssuedAt: token.CreatedAt.UTC().Format(time.RFC3339),
		ExpiresAt: func() *string {
			if token.ExpiresAt != nil {
				expiresAt := token.ExpiresAt.UTC().Format(time.RFC3339)
				return &expiresAt
			}
			return nil
		}(),
	})
}
