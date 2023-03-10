package hook

import (
	"context"
	"github.com/oechsler-it/identity/cqrs"
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
	Model  *model.InMemoryUserRepo
	Create cqrs.CommandHandler[command.Create]
}

func UseCreateRootUser(hook *CreateRootUser) {
	hook.OnStart(hook.onStart)
}

func (e *CreateRootUser) onStart(ctx context.Context) error {
	cmd := command.Create{
		Profile: domain.Profile{
			FirstName: e.Env.String(
				"INITIAL_USER_FIRST_NAME",
				"Initial",
			),
			LastName: e.Env.String(
				"INITIAL_USER_LAST_NAME",
				"User",
			),
		},
		Password: domain.PlainPassword(e.Env.String(
			"INITIAL_USER_PASSWORD",
			"change-me",
		)),
	}

	id, err := e.Model.NextId(ctx)
	if err == nil {
		cmd.Id = id
	} else {
		return err
	}

	if err := e.Create.Handle(ctx, cmd); err != nil {
		return err
	}

	e.Logger.WithField("id", uuid.UUID(id).String()).
		WithField("password", cmd.Password).
		Info("Root user created")

	return nil
}
