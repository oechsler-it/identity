package fiber

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/cqrs"
	sessionFiber "github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
	"github.com/oechsler-it/identity/modules/token/infra/model"
	userDomain "github.com/oechsler-it/identity/modules/user/domain"
	userFiber "github.com/oechsler-it/identity/modules/user/infra/fiber"
	"github.com/oechsler-it/identity/runtime"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type issueTokenRequest struct {
	Description       string   `json:"description" validate:"required"`
	Permissions       []string `json:"permissions"`
	LifetimeInSeconds int      `json:"lifetime_in_seconds" validate:"omitempty"`
}

type issueTokenResponse struct {
	Token     string  `json:"token"`
	Type      string  `json:"type"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

type IssueTokenHandler struct {
	*fiber.App
	// ---
	Logger   *logrus.Logger
	Validate *validator.Validate
	Env      *runtime.Env
	// ---
	RenewMiddleware      *sessionFiber.RenewMiddleware
	ProtectMiddleware    *sessionFiber.ProtectSessionMiddleware
	UserMiddleware       *userFiber.UserMiddleware
	PermissionMiddleware *userFiber.UserPermissionMiddleware
	// ---
	Repo  *model.GormTokenRepo
	Issue cqrs.CommandHandler[command.Issue]
}

func UseIssueTokenHandler(handler *IssueTokenHandler) {
	token := handler.Group("/token")
	token.Post("/",
		handler.RenewMiddleware.Handle,
		handler.ProtectMiddleware.Handle,
		handler.UserMiddleware.Handle,
		handler.PermissionMiddleware.Has("all:token:issue"),
		handler.post)
}

// @Summary	Issue a new token
// @Accept		json
// @Produce	json
// @Param		body	body		issueTokenRequest	true	"Information of the token to issue"
// @Success	201		{object}	issueTokenResponse
// @Failure	400
// @Failure	401
// @Failure	403
// @Failure	422
// @Failure	500
// @Router		/token [post]
// @Tags		Token
func (e *IssueTokenHandler) post(ctx *fiber.Ctx) error {
	var body issueTokenRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	if err := e.Validate.Struct(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, ok := ctx.Locals("user").(*userDomain.User)
	if !ok {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	lifetimeInSeconds := body.LifetimeInSeconds
	if lifetimeInSeconds == 0 {
		lifetimeInSeconds = e.Env.Int("TOKEN_LIFETIME_IN_DAYS") * 60 * 60 * 24
	}

	id, err := e.Repo.NextId(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if err := e.Issue.Handle(ctx.Context(), command.Issue{
		Id:          id,
		Description: body.Description,
		UserId:      domain.UserId(user.Id),
		UserPermissions: lo.Map(user.Permissions, func(permission userDomain.Permission, _ int) domain.Permission {
			return domain.Permission(permission)
		}),
		IncludedPermissions: lo.Map(body.Permissions, func(permission string, _ int) domain.Permission {
			return domain.Permission(permission)
		}),
		LifetimeInSeconds: lifetimeInSeconds,
	}); err != nil {
		if errors.Is(err, domain.ErrTokenCanNotBeGrantedPermission) {
			return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return err
	}

	ctx.Location(fmt.Sprintf("%s/token/%s",
		ctx.BaseURL(),
		id,
	))

	var expiresAt time.Time
	if lifetimeInSeconds > 0 {
		time.Now().Add(time.Duration(lifetimeInSeconds) * time.Second)
	}
	var expiresAtString *string
	if !expiresAt.IsZero() {
		expiresAtValue := expiresAt.Format(time.RFC3339)
		expiresAtString = &expiresAtValue
	}

	e.Logger.WithField("token_id", string(id[:8])).
		WithField("user_id", uuid.UUID(user.Id).String()).
		Info("Token issued")

	return ctx.Status(fiber.StatusCreated).JSON(issueTokenResponse{
		Token:     string(id),
		Type:      "Bearer",
		ExpiresAt: expiresAtString,
	})
}
