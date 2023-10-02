package fiber

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
	"github.com/oechsler-it/identity/modules/user/infra/model"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type createUserRequest struct {
	Password string `json:"password" validate:"gte=8"`
}

type CreateUserHandler struct {
	*fiber.App
	// ---
	Logger   *logrus.Logger
	Validate *validator.Validate
	// ---
	RenewMiddleware      *sessionFiber.RenewMiddleware
	ProtectMiddleware    *sessionFiber.ProtectSessionMiddleware
	UserMiddleware       *UserMiddleware
	PermissionMiddleware *UserPermissionMiddleware
	// ---
	Repo   *model.GormUserRepo
	Create cqrs.CommandHandler[command.Create]
}

func UseCreateUserHandler(handler *CreateUserHandler) {
	user := handler.Group("/user")
	user.Post("/",
		handler.RenewMiddleware.Handle,
		handler.ProtectMiddleware.Handle,
		handler.UserMiddleware.Handle,
		handler.PermissionMiddleware.Has("all:user:create"),
		handler.post)
}

// @Summary	Create a new user
// @Accept		json
// @Produce	text/plain
// @Param		body	body	createUserRequest	true	"Information of the user to create"
// @Success	201
// @Failure	400
// @Failure	401
// @Failure	403
// @Failure	422
// @Failure	500
// @Router		/user [post]
// @Tags		User
func (e *CreateUserHandler) post(ctx *fiber.Ctx) error {
	var body createUserRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	// TODO: Improve password policy (e.g. length, complexity, ...)
	if err := e.Validate.Struct(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			SendString("Password must be at least 8 characters long")
	}

	id, err := e.Repo.NextId(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if err := e.Create.Handle(ctx.Context(), command.Create{
		Id:       id,
		Password: domain.PlainPassword(body.Password),
	}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	e.Logger.WithField("id", uuid.UUID(id).String()).
		Info("User created")

	ctx.Location(fmt.Sprintf("%s/user/%s",
		ctx.BaseURL(),
		uuid.UUID(id).String(),
	))

	return ctx.SendStatus(fiber.StatusCreated)
}
