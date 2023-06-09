package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type sessionOwner struct {
	DeviceId string `json:"device_id"`
	UserId   string `json:"user_id"`
}

type sessionResponse struct {
	Id        string       `json:"id"`
	OwnedBy   sessionOwner `json:"owned_by"`
	ExpiresAt string       `json:"expires_at"`
}

type ActiveSessionHandler struct {
	*fiber.App
	// ---
	RenewMiddleware   *RenewMiddleware
	ProtectMiddleware *ProtectMiddleware
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Session]
}

func UseActiveSessionHandler(handler *ActiveSessionHandler) {
	session := handler.Group("/session")
	session.Use(handler.RenewMiddleware.Handle)
	session.Use(handler.ProtectMiddleware.Handle)
	session.Get("/active", handler.get)
}

//	@Summary	Get details of the active session
//	@Produce	json
//	@Success	200	{object}	sessionResponse
//	@Failure	401
//	@Router		/session/active [get]
//	@Tags		Session
func (e *ActiveSessionHandler) get(ctx *fiber.Ctx) error {
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

	// ---

	return ctx.JSON(sessionResponse{
		Id: uuid.UUID(session.Id).String(),
		OwnedBy: sessionOwner{
			DeviceId: uuid.UUID(session.OwnedBy.DeviceId).String(),
			UserId:   uuid.UUID(session.OwnedBy.UserId).String(),
		},
		ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339),
	})
}
