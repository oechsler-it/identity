package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	sessionDomain "github.com/oechsler-it/identity/modules/session/domain"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/token/app/query"
	"github.com/oechsler-it/identity/modules/token/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type tokenOwner struct {
	UserId string `json:"user_id"`
}

type tokenResponse struct {
	Id          string     `json:"id"`
	Description string     `json:"description"`
	OwnedBy     tokenOwner `json:"owned_by"`
	IssuedAt    string     `json:"issued_at"`
	ExpiresAt   *string    `json:"expires_at,omitempty"`
}

type tokenListResponse []tokenResponse

type ActiveTokensHandler struct {
	*fiber.App
	// ---
	RenewMiddleware          *sessionFiber.RenewMiddleware
	ProtectSessionMiddleware *sessionFiber.ProtectSessionMiddleware
	// ---
	FindByOwnerUserId cqrs.QueryHandler[query.FindByOwnerUserId, []*domain.Token]
}

func UseActiveTokensHandler(handler *ActiveTokensHandler) {
	token := handler.Group("/token")
	token.Get("/",
		handler.RenewMiddleware.Handle,
		handler.ProtectSessionMiddleware.Handle,
		handler.get)
}

// @Summary	List all active tokens belonging to the owner of the current session
// @Produce	json
// @Success	200	{object}	tokenListResponse
// @Failure	401
// @Failure	500
// @Router		/token [get]
// @Tags		Token
func (e *ActiveTokensHandler) get(ctx *fiber.Ctx) error {
	session, ok := ctx.Locals("session").(*sessionDomain.Session)
	if !ok {
		return fiber.ErrInternalServerError
	}

	tokens, err := e.FindByOwnerUserId.Handle(ctx.Context(), query.FindByOwnerUserId{
		UserId: domain.UserId(session.OwnedBy.UserId),
	})
	if err != nil {
		return err
	}

	response := make(tokenListResponse, len(tokens))
	for i, token := range tokens {
		response[i] = tokenResponse{
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
		}
	}

	return ctx.JSON(response)
}
