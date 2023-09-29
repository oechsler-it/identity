package permission

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/cqrs"
	commandHandler "github.com/oechsler-it/identity/modules/permission/app"
	queryHandler "github.com/oechsler-it/identity/modules/permission/app"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/app/query"
	"github.com/oechsler-it/identity/modules/permission/domain"
	"github.com/oechsler-it/identity/modules/permission/infra/fiber"
	"github.com/oechsler-it/identity/modules/permission/infra/hook"
	"github.com/oechsler-it/identity/modules/permission/infra/model"
)

type Options struct {
	CreateBasePermissions *hook.CreateBasePermissions
	CreateHandler         *fiber.CreateHandler
	DeleteHandler         *fiber.DeleteHandler
	PermissionsHandler    *fiber.PermissionsHandler
}

func UsePermission(opts *Options) {
	hook.UseCreateBasePermissions(opts.CreateBasePermissions)
	fiber.UseCreateHandler(opts.CreateHandler)
	fiber.UseDeleteHandler(opts.DeleteHandler)
	fiber.UsePermissionsHandler(opts.PermissionsHandler)
}

var WirePermission = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewCreateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Create]), new(*commandHandler.CreateHandler)),

	commandHandler.NewDeleteHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Delete]), new(*commandHandler.DeleteHandler)),

	commandHandler.NewVerifyPermissionNotExistsHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyPermissionNotExists]), new(*commandHandler.VerifyPermissionNotExistsHandler)),

	queryHandler.NewFindAllHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindAll, []*domain.Permission]), new(*queryHandler.FindAllHandler)),

	queryHandler.NewFindByNameHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindByName, *domain.Permission]), new(*queryHandler.FindByNameHandler)),

	wire.Struct(new(fiber.CreateHandler), "*"),
	wire.Struct(new(fiber.DeleteHandler), "*"),
	wire.Struct(new(fiber.PermissionsHandler), "*"),
	wire.Struct(new(hook.CreateBasePermissions), "*"),

	model.NewGormPermissionRepo,
	wire.Bind(new(commandHandler.CreateWriteModel), new(*model.GormPermissionRepo)),
	wire.Bind(new(commandHandler.DeleteWriteModel), new(*model.GormPermissionRepo)),
	wire.Bind(new(commandHandler.VerifyPermissionNotExistsReadModel), new(*model.GormPermissionRepo)),
	wire.Bind(new(queryHandler.FindAllReadModel), new(*model.GormPermissionRepo)),
	wire.Bind(new(queryHandler.FindByNameReadModel), new(*model.GormPermissionRepo)),
)
