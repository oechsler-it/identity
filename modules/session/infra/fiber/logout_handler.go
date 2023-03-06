package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
)

type LogoutHandler struct {
	*fiber.App
	// ---
	ProtectMiddleware *ProtectMiddleware
	// ---
	Revoke cqrs.CommandHandler[command.Revoke]
}

func UseLogoutHandler(handler *LogoutHandler) {
	logout := handler.Group("/logout")
	logout.Use(handler.ProtectMiddleware.Handle)
	logout.Post("/", handler.post)
}

func (e *LogoutHandler) post(ctx *fiber.Ctx) error {
	sessionIdCookie := ctx.Cookies("session_id")

	sessionId, err := uuid.FromString(sessionIdCookie)
	if err != nil {
		return err
	}

	if err = e.Revoke.Handle(ctx.Context(), command.Revoke{
		Id: domain.SessionId(sessionId),
	}); err != nil {
		return err
	}

	ctx.ClearCookie("session_id")

	return ctx.SendStatus(fiber.StatusOK)
}
