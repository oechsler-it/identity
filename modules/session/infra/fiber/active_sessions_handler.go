package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type sessionsHandlerResponses []sessionResponse

type ActiveSessionsHandler struct {
	*fiber.App
	// ---
	RenewMiddleware   *RenewMiddleware
	ProtectMiddleware *ProtectMiddleware
	// ---
	FindById          cqrs.QueryHandler[query.FindById, *domain.Session]
	FindByOwnerUserId cqrs.QueryHandler[query.FindByOwnerUserId, []*domain.Session]
}

func UseActiveSessionsHandler(handler *ActiveSessionsHandler) {
	session := handler.Group("/session")
	session.Use(handler.RenewMiddleware.Handle)
	session.Use(handler.ProtectMiddleware.Handle)
	session.Get("/", handler.get)
}

//	@Summary	List all active sessions belonging to the owner of the current session
//	@Produce	json
//	@Success	200	{object}	sessionsHandlerResponses
//	@Failure	401
//	@Router		/session [get]
//	@Tags		Session
func (e *ActiveSessionsHandler) get(ctx *fiber.Ctx) error {
	sessionIdCookie := ctx.Cookies("session_id")

	sessionId, err := uuid.FromString(sessionIdCookie)
	if err != nil {
		return err
	}

	session, err := e.FindById.Handle(ctx.Context(), query.FindById{
		Id: domain.SessionId(sessionId),
	})
	if err != nil {
		return err
	}

	sessions, err := e.FindByOwnerUserId.Handle(ctx.Context(), query.FindByOwnerUserId{
		UserId: session.OwnedBy.UserId,
	})
	if err != nil {
		return err
	}

	// ---

	response := make(sessionsHandlerResponses, len(sessions))
	for i, session := range sessions {
		response[i] = sessionResponse{
			Id: uuid.UUID(session.Id).String(),
			OwnedBy: sessionOwner{
				DeviceId: uuid.UUID(session.OwnedBy.DeviceId).String(),
				UserId:   uuid.UUID(session.OwnedBy.UserId).String(),
			},
			ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339),
		}
	}
	return ctx.JSON(response)
}
