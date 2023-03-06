package fiber

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	"github.com/oechsler-it/identity/modules/session/infra/model"
	userCommand "github.com/oechsler-it/identity/modules/user/app/command"
	userQuery "github.com/oechsler-it/identity/modules/user/app/query"
	userDomain "github.com/oechsler-it/identity/modules/user/domain"
	"github.com/oechsler-it/identity/runtime"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type LoginHandler struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
	Env    *runtime.Env
	// ---
	Model    *model.InMemorySessionRepo
	Initiate cqrs.CommandHandler[command.Initiate]
	// ---
	FindUserByIdentifier cqrs.QueryHandler[userQuery.FindByIdentifier, *userDomain.User]
	VerifyPassword       cqrs.CommandHandler[userCommand.VerifyPassword]
}

func UseFiberLoginHandler(handler *LoginHandler) {
	handler.Post("/login", handler.handle)
}

func (e *LoginHandler) handle(ctx *fiber.Ctx) error {
	if ctx.Get("Content-Type") != "application/x-www-form-urlencoded" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid content type")
	}

	user, err := e.FindUserByIdentifier.Handle(ctx.Context(), userQuery.FindByIdentifier{
		Identifier: ctx.FormValue("identifier"),
	})
	if err != nil {
		if errors.Is(err, userDomain.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return err
	}

	if err := e.VerifyPassword.Handle(ctx.Context(), userCommand.VerifyPassword{
		Id:       user.Id,
		Password: userDomain.PlainPassword(ctx.FormValue("password")),
	}); err != nil {
		if errors.Is(err, userDomain.ErrInvalidPassword) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return err
	}

	sessionId, err := e.Model.NextId(ctx.Context())
	if err != nil {
		return err
	}

	deviceId, ok := ctx.Locals("device_id").(uuid.UUID)
	if !ok {
		return domain.ErrInvalidDeviceId
	}

	lifetimeInSeconds := e.Env.Int("SESSION_LIFETIME_IN_HOURS", 8) * 60 * 60

	if err := e.Initiate.Handle(ctx.Context(), command.Initiate{
		Id:                sessionId,
		UserId:            domain.UserId(user.Id),
		DeviceId:          domain.DeviceId(deviceId),
		LifetimeInSeconds: lifetimeInSeconds,
		Renewable:         ctx.FormValue("renewable") == "true",
	}); err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Value:   uuid.UUID(sessionId).String(),
		Path:    "/",
		Expires: time.Now().Add(time.Duration(lifetimeInSeconds) * time.Second),
	})

	return ctx.SendStatus(fiber.StatusOK)
}