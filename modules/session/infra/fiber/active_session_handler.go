package fiber

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	uuid "github.com/satori/go.uuid"
)

type sessionOwner struct {
	DeviceId string `json:"device_id"`
	UserId   string `json:"user_id"`
}

type sessionResponse struct {
	Id          string       `json:"id"`
	OwnedBy     sessionOwner `json:"owned_by"`
	Active      bool         `json:"active"`
	InitiatedAt string       `json:"initiated_at"`
	ExpiresAt   string       `json:"expires_at"`
}

type ActiveSessionHandler struct {
	*fiber.App
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Session]
}

func UseActiveSessionHandler(handler *ActiveSessionHandler) {
	session := handler.Group("/session")
	session.Get("/active",
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.AuthenticatedMiddleware.Handle,
		// ---
		handler.get)
}

// @Summary	Get details of the active session
// @Produce	json
// @Success	200	{object}	sessionResponse
// @Failure	401
// @Failure	500
// @Router		/session/active [get]
// @Tags		Session
func (e *ActiveSessionHandler) get(ctx *fiber.Ctx) error {
	session, ok := ctx.Locals("session").(*domain.Session)
	if !ok {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(sessionResponse{
		Id: uuid.UUID(session.Id).String(),
		OwnedBy: sessionOwner{
			DeviceId: uuid.UUID(session.OwnedBy.DeviceId).String(),
			UserId:   uuid.UUID(session.OwnedBy.UserId).String(),
		},
		Active:      true,
		InitiatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
		ExpiresAt:   session.ExpiresAt.UTC().Format(time.RFC3339),
	})
}
