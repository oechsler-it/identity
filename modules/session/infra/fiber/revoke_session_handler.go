package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type RevokeSessionHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	ProtectMiddleware *ProtectMiddleware
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Session]
	Revoke   cqrs.CommandHandler[command.Revoke]
}

func UseRevokeSessionHandler(handler *RevokeSessionHandler) {
	session := handler.Group("/session")
	session.Use(handler.ProtectMiddleware.Handle)
	session.Delete("/revoke/:id", handler.delete)
}

//	@Summary	Revoke a session
//	@Produce	text/plain
//	@Param		id	path	string	true	"Id of the session"
//	@Success	204
//	@Failure	401
//	@Failure	403
//	@Failure	404
//	@Router		/session/revoke/{id} [delete]
//	@Tags		Session
func (e *RevokeSessionHandler) delete(ctx *fiber.Ctx) error {
	sessionIdCookie := ctx.Cookies("session_id")

	sessionId, err := uuid.FromString(sessionIdCookie)
	if err != nil {
		return err
	}

	activeSession, err := e.FindById.Handle(ctx.Context(), query.FindById{
		Id: domain.SessionId(sessionId),
	})
	if err != nil {
		return err
	}

	sessionIdParam := ctx.Params("id")

	sessionId, err = uuid.FromString(sessionIdParam)
	if err != nil {
		return err
	}

	revokeSession, err := e.FindById.Handle(ctx.Context(), query.FindById{
		Id: domain.SessionId(sessionId),
	})
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return err
	}

	// ---

	if err := e.Revoke.Handle(ctx.Context(), command.Revoke{
		Id:             revokeSession.Id,
		RevokingEntity: activeSession.OwnedBy,
	}); err != nil {
		if errors.Is(err, domain.ErrSessionIsExpired) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithFields(logrus.Fields{
		"session_id": revokeSession.Id,
	}).Info("Session revoked")

	return ctx.SendStatus(fiber.StatusNoContent)
}