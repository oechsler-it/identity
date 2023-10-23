package fiber

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
	uuid "github.com/satori/go.uuid"
)

type SessionByIdHandler struct {
	*fiber.App
	// ---
	TokenAuthMiddleware *tokenFiberMiddleware.TokenAuthMiddleware
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Session]
}

func UseSessionByIdHandler(handler *SessionByIdHandler) {
	session := handler.Group("/session")
	session.Get("/:id",
		handler.TokenAuthMiddleware.Handle,
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.AuthenticatedMiddleware.Handle,
		// ---
		handler.get)
}

//	@Summary	Get details of a session
//	@Produce	json
//	@Param		id	path		string	true	"Session Id"
//	@Success	200	{object}	sessionResponse
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/session/{id} [get]
//	@Security	TokenAuth
//	@Tags		Session
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
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return err
	}

	return ctx.JSON(sessionResponse{
		Id: uuid.UUID(session.Id).String(),
		OwnedBy: sessionOwner{
			DeviceId: uuid.UUID(session.OwnedBy.DeviceId).String(),
			UserId:   uuid.UUID(session.OwnedBy.UserId).String(),
		},
		ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339),
	})
}
