package hook

import (
	"context"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
	"github.com/oechsler-it/identity/runtime"
	"github.com/sirupsen/logrus"
)

type CreateBasePermissions struct {
	*runtime.Hooks
	// ---
	Logger *logrus.Logger
	Env    *runtime.Env
	// ---
	VerifyPermissionNotExists cqrs.CommandHandler[command.VerifyPermissionNotExists]
	Create                    cqrs.CommandHandler[command.Create]
}

func UseCreateBasePermissions(hook *CreateBasePermissions) {
	hook.OnStart(hook.all)
	hook.OnStart(hook.allPermission)
	hook.OnStart(hook.allPermissionCreate)
	hook.OnStart(hook.allPermissionDelete)
	hook.OnStart(hook.allUser)
	hook.OnStart(hook.allUserCreate)
	hook.OnStart(hook.allUserDelete)
	hook.OnStart(hook.allUserPermission)
	hook.OnStart(hook.allUserPermissionGrant)
	hook.OnStart(hook.allUserPermissionRevoke)
}

func (e *CreateBasePermissions) ensureCreated(ctx context.Context, name string, description string) error {
	if err := e.VerifyPermissionNotExists.Handle(ctx, command.VerifyPermissionNotExists{
		Name: domain.PermissionName(name),
	}); err != nil {
		return nil
	}

	if err := e.Create.Handle(ctx, command.Create{
		Name:        domain.PermissionName(name),
		Description: description,
	}); err != nil {
		return err
	}

	e.Logger.WithField("name", name).
		Info("Permission created")

	return nil
}

func (e *CreateBasePermissions) all(ctx context.Context) error {
	return e.ensureCreated(ctx, "all", "Root permission")
}

func (e *CreateBasePermissions) allPermission(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:permission", "Manage permissions")
}

func (e *CreateBasePermissions) allPermissionCreate(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:permission:create", "Create new permissions")
}

func (e *CreateBasePermissions) allPermissionDelete(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:permission:delete", "Delete permissions")
}

func (e *CreateBasePermissions) allUser(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:user", "Manage users")
}

func (e *CreateBasePermissions) allUserCreate(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:user:create", "Create new users")
}

func (e *CreateBasePermissions) allUserDelete(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:user:delete", "Delete users")
}

func (e *CreateBasePermissions) allUserPermission(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:user:permission", "Manage users permissions")
}

func (e *CreateBasePermissions) allUserPermissionGrant(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:user:permission:grant", "Grant permissions to users")
}

func (e *CreateBasePermissions) allUserPermissionRevoke(ctx context.Context) error {
	return e.ensureCreated(ctx, "all:user:permission:revoke", "Revoke permissions from users")
}
