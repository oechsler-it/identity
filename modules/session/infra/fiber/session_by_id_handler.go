package fiber

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type SessionByIdHandler struct {
	*fiber.App
	// ---
	RenewMiddleware   *RenewMiddleware
	ProtectMiddleware *ProtectMiddleware
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Session]
}

func UseSessionByIdHandler(handler *SessionByIdHandler) {
	session := handler.Group("/session")
	session.Use(handler.RenewMiddleware.Handle)
	session.Use(handler.ProtectMiddleware.Handle)
	session.Get("/:id", handler.get)
}

// @Summary	Get details of a session
// @Produce	json
// @Param		id	path		string	true	"Session Id"
// @Success	200	{object}	sessionResponse
// @Failure	401
// @Failure	404
// @Failure	500
// @Router		/session/{id} [get]
// @Tags		Session
func (e *SessionByIdHandler) get(ctx *fiber.Ctx) error {
	sessionId, err := uuid.FromString(ctx.Params("id"))
	if err != nil {
		return err
	}

	session, err := e.FindById.Handle(ctx.Context(), query.FindById{
		Id: domain.SessionId(sessionId),
	})
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
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
