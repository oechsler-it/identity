package fiber

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
	"github.com/oechsler-it/identity/runtime"
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
	Env    *runtime.Env
	// ---
	Create cqrs.CommandHandler[command.Create]
}

func UseCreateHandler(handler *CreateHandler) {
	create := handler.Group("/permission")
	create.Post("/", handler.post)
}

//	@Summary	Create a new permission
//	@Accept		json
//	@Produce	text/plain
//	@Param		command	body	createPermissionRequest	true	"Create command"
//	@Success	201
//	@Failure	400
//	@Router		/permission [post]
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

	return ctx.SendStatus(fiber.StatusCreated)
}
