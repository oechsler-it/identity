package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	middlewareFiber "github.com/oechsler-it/identity/modules/middleware/infra/fiber"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
	sessionFiberMiddleware "github.com/oechsler-it/identity/modules/session/infra/fiber/middleware"
	tokenFiberMiddleware "github.com/oechsler-it/identity/modules/token/infra/fiber/middleware"
	userFiberMiddleware "github.com/oechsler-it/identity/modules/user/infra/fiber/middleware"
	"github.com/sirupsen/logrus"
)

type createPermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateHandler struct {
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
	Create cqrs.CommandHandler[command.Create]
}

func UseCreateHandler(handler *CreateHandler) {
	create := handler.Group("/permission")
	create.Post("/",
		handler.TokenAuthMiddleware.Handle,
		handler.TokenPermissionMiddleware.Has("all:permission:create"),
		// ---
		handler.RenewMiddleware.Handle,
		handler.SessionAuthMiddleware.Handle,
		// ---
		handler.UserMiddleware.Handle,
		handler.UserPermissionMiddleware.Has("all:permission:create"),
		// ---
		handler.AuthenticatedMiddleware.Handle,
		handler.AuthorizedMiddleware.Handle,
		// ---
		handler.post)
}

//	@Summary	Create a new permission
//	@Accept		json
//	@Produce	text/plain
//	@Param		command	body	createPermissionRequest	true	"Create command"
//	@Success	201
//	@Failure	400
//	@Failure	401
//	@Failure	500
//	@Router		/permission [post]
//	@Security	TokenAuth
//	@Tags		Permission
func (e *CreateHandler) post(ctx *fiber.Ctx) error {
	if ctx.Get("Content-Type") != "application/json" {
		return fiber.ErrUnsupportedMediaType
	}

	var dto createPermissionRequest
	if err := ctx.BodyParser(&dto); err != nil {
		return err
	}

	if err := e.Create.Handle(ctx.Context(), command.Create{
		Name:        domain.PermissionName(dto.Name),
		Description: dto.Description,
	}); err != nil {
		if errors.Is(err, domain.ErrPermissionAlreadyExists) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return err
	}

	e.Logger.WithFields(logrus.Fields{
		"name": dto.Name,
	}).Info("Permission created")

	ctx.Set("Location", "/permission/"+dto.Name)

	return ctx.SendStatus(fiber.StatusCreated)
}
