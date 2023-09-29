package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	userFiber "github.com/oechsler-it/identity/modules/user/infra/fiber"
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
	RenewMiddleware      *sessionFiber.RenewMiddleware
	ProtectMiddleware    *sessionFiber.ProtectMiddleware
	UserMiddleware       *userFiber.UserMiddleware
	PermissionMiddleware *userFiber.PermissionMiddleware
	// ---
	Create cqrs.CommandHandler[command.Create]
}

func UseCreateHandler(handler *CreateHandler) {
	create := handler.Group("/permission")
	create.Post("/",
		handler.RenewMiddleware.Handle,
		handler.ProtectMiddleware.Handle,
		handler.UserMiddleware.Handle,
		handler.PermissionMiddleware.Has("all:permission:create"),
		handler.post)
}

// @Summary	Create a new permission
// @Accept		json
// @Produce	text/plain
// @Param		command	body	createPermissionRequest	true	"Create command"
// @Success	201
// @Failure	400
// @Failure	401
// @Failure	500
// @Router		/permission [post]
// @Tags		Permission
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
