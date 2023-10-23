package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	sessionDomain "github.com/oechsler-it/identity/modules/session/domain"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
	userDomain "github.com/oechsler-it/identity/modules/user/domain"
	userFiberMiddleware "github.com/oechsler-it/identity/modules/user/infra/fiber/middleware"
	"github.com/sirupsen/logrus"
)

type RevokeTokenHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	TokenAuthMiddleware       *tokenFiberMiddleware.TokenAuthMiddleware
	TokenPermissionMiddleware *tokenFiberMiddleware.TokenPermissionMiddleware
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	UserMiddleware           *userFiberMiddleware.UserMiddleware
	UserPermissionMiddleware *userFiberMiddleware.UserPermissionMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	AuthorizedMiddleware    *middlewareFiber.AuthorizedMiddleware
	// ---
	Revoke cqrs.CommandHandler[command.Revoke]
}

func UseRevokeTokenHandler(handler *RevokeTokenHandler) {
	token := handler.Group("/token")
	token.Delete("/revoke/:id",
		handler.TokenAuthMiddleware.Handle,
		handler.TokenPermissionMiddleware.Has("all:token:revoke"),
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.UserMiddleware.Handle,
		handler.UserPermissionMiddleware.Has(userDomain.PermissionNone),
		// ---
		handler.AuthenticatedMiddleware.Handle,
		handler.AuthorizedMiddleware.Handle,
		// ---
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
// @Security	TokenAuth
// @Tags		Token
func (e *RevokeTokenHandler) delete(ctx *fiber.Ctx) error {
	token, tokenOk := ctx.Locals("token").(*domain.Token)
	activeSession, activeSessionOk := ctx.Locals("session").(*sessionDomain.Session)
	if !tokenOk && !activeSessionOk {
		return fiber.ErrInternalServerError
	}

	idParam := ctx.Params("id")

	idPartial, err := domain.NewTokenIdPartialFromString(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid token id")
	}

	revokingEntity := func() domain.Owner {
		if token != nil {
			return token.OwnedBy
		}
		return domain.Owner{
			UserId: domain.UserId(activeSession.OwnedBy.UserId),
		}
	}()

	if err := e.Revoke.Handle(ctx.Context(), command.Revoke{
		IdPartial:      idPartial,
		RevokingEntity: revokingEntity,
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
