package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type LogoutHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	ProtectMiddleware *ProtectMiddleware
	// ---
	Revoke cqrs.CommandHandler[command.Revoke]
}

func UseLogoutHandler(handler *LogoutHandler) {
	logout := handler.Group("/logout")
	logout.Use(handler.ProtectMiddleware.Handle)
	logout.Delete("/", handler.delete)
}

//	@Summary	Revoke the current session
//	@Produce	text/plain
//	@Success	204
//	@Failure	401
//	@Router		/logout [delete]
//	@Tags		Session
func (e *LogoutHandler) delete(ctx *fiber.Ctx) error {
	sessionId := ctx.Locals("session_id").(domain.SessionId)

	if err := e.Revoke.Handle(ctx.Context(), command.Revoke{
		Id: sessionId,
	}); err != nil {
		return err
	}

	ctx.ClearCookie("session_id")

	e.Logger.WithFields(logrus.Fields{
		"session_id": uuid.UUID(sessionId).String(),
	}).Info("session revoked")

	return ctx.SendStatus(fiber.StatusNoContent)
}
