package fiber

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	sessionDomain "github.com/oechsler-it/identity/modules/session/domain"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
	"github.com/sirupsen/logrus"
)

type RevokeTokenHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	RenewMiddleware          *sessionFiber.RenewMiddleware
	ProtectSessionMiddleware *sessionFiber.ProtectSessionMiddleware
	// ---
	Revoke cqrs.CommandHandler[command.Revoke]
}

func UseRevokeTokenHandler(handler *RevokeTokenHandler) {
	token := handler.Group("/token")
	token.Delete("/revoke/:id",
		handler.RenewMiddleware.Handle,
		handler.ProtectSessionMiddleware.Handle,
		handler.delete)
}

// @Summary	Revoke a token
// @Produce	text/plain
// @Param		id	path	string	true	"Id of the token"
// @Success	204
// @Failure	400
// @Failure	401
// @Failure	403
// @Failure	404
// @Failure	500
// @Router		/token/revoke/{id} [delete]
// @Tags		Token
func (e *RevokeTokenHandler) delete(ctx *fiber.Ctx) error {
	activeSession, ok := ctx.Locals("session").(*sessionDomain.Session)
	if !ok {
		return fiber.ErrInternalServerError
	}

	idParam := ctx.Params("id")

	idPartial, err := domain.NewTokenIdPartialFromString(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid token id")
	}

	if err := e.Revoke.Handle(ctx.Context(), command.Revoke{
		IdPartial: idPartial,
		RevokingEntity: domain.Owner{
			UserId: domain.UserId(activeSession.OwnedBy.UserId),
		},
	}); err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrTokenDoesNotBelongToOwner) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		if errors.Is(err, domain.ErrTokenIsExpired) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithFields(logrus.Fields{
		"token_id": idPartial.String(),
	}).Info("token revoked")

	return ctx.SendStatus(fiber.StatusNoContent)
}
