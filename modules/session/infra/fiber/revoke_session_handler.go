package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	tokenDomain "github.com/oechsler-it/identity/modules/token/domain"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
	userDomain "github.com/oechsler-it/identity/modules/user/domain"
	userFiberMiddleware "github.com/oechsler-it/identity/modules/user/infra/fiber/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type RevokeSessionHandler struct {
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

func UseRevokeSessionHandler(handler *RevokeSessionHandler) {
	session := handler.Group("/session")
	session.Delete("/revoke/:id",
		handler.TokenAuthMiddleware.Handle,
		handler.TokenPermissionMiddleware.Has("all:session:revoke"),
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

//	@Summary	Revoke a session belonging to the owner of the current session
//	@Produce	text/plain
//	@Param		id	path	string	true	"Id of the session"
//	@Success	204
//	@Failure	400
//	@Failure	401
//	@Failure	403
//	@Failure	404
//	@Failure	500
//	@Router		/session/revoke/{id} [delete]
//	@Security	TokenAuth
//	@Tags		Session
func (e *RevokeSessionHandler) delete(ctx *fiber.Ctx) error {
	token, tokenOk := ctx.Locals("token").(*tokenDomain.Token)
	activeSession, activeSessionOk := ctx.Locals("session").(*domain.Session)
	if !tokenOk && !activeSessionOk {
		return fiber.ErrInternalServerError
	}

	revokeSessionIdParam := ctx.Params("id")

	revokeSessionId, err := uuid.FromString(revokeSessionIdParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	revokingEntity := func() domain.Owner {
		if token != nil {
			return domain.Owner{
				UserId: domain.UserId(token.OwnedBy.UserId),
			}
		}
		return activeSession.OwnedBy
	}()

	if err := e.Revoke.Handle(ctx.Context(), command.Revoke{
		Id:             domain.SessionId(revokeSessionId),
		RevokingEntity: revokingEntity,
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
