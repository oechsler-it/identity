package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type LogoutHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	// ---
	FindById cqrs.QueryHandler[query.FindById, *domain.Session]
	Revoke   cqrs.CommandHandler[command.Revoke]
}

func UseLogoutHandler(handler *LogoutHandler) {
	logout := handler.Group("/logout")
	logout.Delete("/",
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.AuthenticatedMiddleware.Handle,
		// ---
		handler.delete)
}

// @Summary	Revoke the current session
// @Produce	text/plain
// @Success	204
// @Failure	401
// @Failure	500
// @Router		/logout [delete]
// @Tags		Session
func (e *LogoutHandler) delete(ctx *fiber.Ctx) error {
	session, ok := ctx.Locals("session").(*domain.Session)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := e.Revoke.Handle(ctx.Context(), command.Revoke{
		Id:             session.Id,
		RevokingEntity: session.OwnedBy,
	}); err != nil {
		return err
	}

	ctx.ClearCookie("session_id")

	e.Logger.WithFields(logrus.Fields{
		"session_id": uuid.UUID(session.Id).String(),
	}).Info("Session revoked")

	return ctx.SendStatus(fiber.StatusNoContent)
}
