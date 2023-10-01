package hook

import (
	"context"
	"github.com/oechsler-it/identity/cqrs"
	permissionQuery "github.com/oechsler-it/identity/modules/permission/app/query"
	permissionDomain "github.com/oechsler-it/identity/modules/permission/domain"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
	"github.com/oechsler-it/identity/modules/user/infra/model"
	"github.com/oechsler-it/identity/runtime"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type CreateRootUser struct {
	*runtime.Hooks
	// ---
	Logger *logrus.Logger
	Env    *runtime.Env
	// ---
	Repo                 *model.GormUserRepo
	VerifyNoUserExists   cqrs.CommandHandler[command.VerifyNoUserExists]
	Create               cqrs.CommandHandler[command.Create]
	FindPermissionByName cqrs.QueryHandler[permissionQuery.FindByName, *permissionDomain.Permission]
	Grant                cqrs.CommandHandler[command.GrantPermission]
}

func UseCreateRootUser(hook *CreateRootUser) {
	hook.OnStart(hook.onStart)
}

func (e *CreateRootUser) onStart(ctx context.Context) error {
	if err := e.VerifyNoUserExists.Handle(ctx, command.VerifyNoUserExists{}); err != nil {
		return nil
	}

	id, err := e.Repo.NextId(ctx)
	if err != nil {
		return err
	}

	password := domain.PlainPassword(e.Env.String(
		"INITIAL_USER_PASSWORD",
		"change-me",
	))

	if err := e.Create.Handle(ctx, command.Create{
		Id:       id,
		Password: password,
	}); err != nil {
		return err
	}

	e.Logger.WithField("id", uuid.UUID(id).String()).
		WithField("password", password).
		Info("Root user created")

	go func() {
		var permission *permissionDomain.Permission
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
				var err error
				permission, err = e.FindPermissionByName.Handle(ctx, permissionQuery.FindByName{
					Name: "all",
				})
				if err == nil {
					break loop
				}
			}
		}

		if err := e.Grant.Handle(ctx, command.GrantPermission{
			Id:         id,
			Permission: domain.Permission(permission.Name),
		}); err != nil {
			e.Logger.WithError(err).
				WithField("id", uuid.UUID(id).String()).
				WithField("permission", permission.Name).
				Error("Failed to grant permission to root user")
		}

		e.Logger.WithField("id", uuid.UUID(id).String()).
			WithField("permission", permission.Name).
			Info("Granted permission to root user")
	}()

	return nil
}
