package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type RevokeSessionHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	RenewMiddleware          *RenewMiddleware
	ProtectSessionMiddleware *ProtectSessionMiddleware
	// ---
	Revoke cqrs.CommandHandler[command.Revoke]
}

func UseRevokeSessionHandler(handler *RevokeSessionHandler) {
	session := handler.Group("/session")
	session.Delete("/revoke/:id",
		handler.RenewMiddleware.Handle,
		handler.ProtectSessionMiddleware.Handle,
		handler.delete)
}

// @Summary	Revoke a session
// @Produce	text/plain
// @Param		id	path	string	true	"Id of the session"
// @Success	204
// @Failure	400
// @Failure	401
// @Failure	403
// @Failure	404
// @Failure	500
// @Router		/session/revoke/{id} [delete]
// @Tags		Session
func (e *RevokeSessionHandler) delete(ctx *fiber.Ctx) error {
	activeSession, ok := ctx.Locals("session").(*domain.Session)
	if !ok {
		return fiber.ErrInternalServerError
	}

	revokeSessionIdParam := ctx.Params("id")

	revokeSessionId, err := uuid.FromString(revokeSessionIdParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := e.Revoke.Handle(ctx.Context(), command.Revoke{
		Id:             domain.SessionId(revokeSessionId),
		RevokingEntity: activeSession.OwnedBy,
	}); err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrSessionDoesNotBelongToOwner) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrSessionIsExpired) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithFields(logrus.Fields{
		"session_id": revokeSessionId.String(),
	}).Info("Session revoked")

	return ctx.SendStatus(fiber.StatusNoContent)
}
