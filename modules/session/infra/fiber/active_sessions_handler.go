package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type sessionListResponse []sessionResponse

type ActiveSessionsHandler struct {
	*fiber.App
	// ---
	RenewMiddleware          *RenewMiddleware
	ProtectSessionMiddleware *ProtectSessionMiddleware
	// ---
	FindByOwnerUserId cqrs.QueryHandler[query.FindByOwnerUserId, []*domain.Session]
}

func UseActiveSessionsHandler(handler *ActiveSessionsHandler) {
	session := handler.Group("/session")
	session.Get("/",
		handler.RenewMiddleware.Handle,
		handler.ProtectSessionMiddleware.Handle,
		handler.get)
}

// @Summary	List all active sessions belonging to the owner of the current session
// @Produce	json
// @Success	200	{object}	sessionListResponse
// @Failure	401
// @Failure	500
// @Router		/session [get]
// @Tags		Session
func (e *ActiveSessionsHandler) get(ctx *fiber.Ctx) error {
	activeSession, ok := ctx.Locals("session").(*domain.Session)
	if !ok {
		return fiber.ErrInternalServerError
	}

	sessions, err := e.FindByOwnerUserId.Handle(ctx.Context(), query.FindByOwnerUserId{
		UserId: activeSession.OwnedBy.UserId,
	})
	if err != nil {
		return err
	}

	response := make(sessionListResponse, len(sessions))
	for i, session := range sessions {
		response[i] = sessionResponse{
			Id: uuid.UUID(session.Id).String(),
			OwnedBy: sessionOwner{
				DeviceId: uuid.UUID(session.OwnedBy.DeviceId).String(),
				UserId:   uuid.UUID(session.OwnedBy.UserId).String(),
			},
			Active:      activeSession.Id == session.Id,
			InitiatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
			ExpiresAt:   session.ExpiresAt.UTC().Format(time.RFC3339),
		}
	}
	return ctx.JSON(response)
}
