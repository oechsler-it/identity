package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
	"github.com/oechsler-it/identity/modules/user/app/query"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type userResponse struct {
	Id           string `json:"id"`
	RegisteredAt string `json:"registered_at"`
}

type UserByIdHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	// ---
	TokenAuthMiddleware *tokenFiberMiddleware.TokenAuthMiddleware
	// ---
	RenewMiddleware       *sessionFiberMiddleware.RenewMiddleware
	SessionAuthMiddleware *sessionFiberMiddleware.SessionAuthMiddleware
	// ---
	AuthenticatedMiddleware *middlewareFiber.AuthenticatedMiddleware
	// ---
	FindByIdentifier cqrs.QueryHandler[query.FindByIdentifier, *domain.User]
}

func UseUserByIdHandler(handler *UserByIdHandler) {
	user := handler.Group("/user")
	user.Get("/:id",
		handler.TokenAuthMiddleware.Handle,
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.AuthenticatedMiddleware.Handle,
		// ---
		handler.get)
}

// @Summary	Get information about a user
// @Produce	json
// @Param		id	path		string	true	"Id of the user"
// @Success	200	{object}	userResponse
// @Failure	400
// @Failure	404
// @Failure	500
// @Router		/user/{id} [get]
// @Security	TokenAuth
// @Tags		User
func (e *UserByIdHandler) get(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := uuid.FromString(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, err := e.FindByIdentifier.Handle(ctx.Context(), query.FindByIdentifier{
		Identifier: id.String(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return err
	}

	return ctx.JSON(userResponse{
		Id:           uuid.UUID(user.Id).String(),
		RegisteredAt: user.CreatedAt.String(),
	})
}
