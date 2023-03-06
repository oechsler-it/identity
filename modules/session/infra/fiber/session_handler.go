package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type sessionHandlerOwner struct {
	DeviceId string `json:"device_id"`
	UserId   string `json:"user_id"`
}

type sessionHandlerResponse struct {
	Id        string              `json:"id"`
	OwnedBy   sessionHandlerOwner `json:"owned_by"`
	ExpiresAt string              `json:"expires_at"`
}

type SessionHandler struct {
	*fiber.App
	// ---
	RenewMiddleware   *RenewMiddleware
	ProtectMiddleware *ProtectMiddleware
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Session]
}

func UseSessionHandler(handler *SessionHandler) {
	session := handler.Group("/session")
	session.Use(handler.RenewMiddleware.Handle)
	session.Use(handler.ProtectMiddleware.Handle)
	session.Get("/", handler.get)
}

func (e *SessionHandler) get(ctx *fiber.Ctx) error {
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

	return ctx.JSON(sessionHandlerResponse{
		Id: uuid.UUID(session.Id).String(),
		OwnedBy: sessionHandlerOwner{
			DeviceId: uuid.UUID(session.OwnedBy.DeviceId).String(),
			UserId:   uuid.UUID(session.OwnedBy.UserId).String(),
		},
		ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339),
	})
}
