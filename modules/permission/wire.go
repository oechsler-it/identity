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
	"github.com/oechsler-it/identity/modules/permission/infra/model"
)

type Options struct {
	CreateHandler      *fiber.CreateHandler
	PermissionsHandler *fiber.PermissionsHandler
}

func UsePermission(opts *Options) {
	fiber.UseCreateHandler(opts.CreateHandler)
	fiber.UsePermissionsHandler(opts.PermissionsHandler)
}

var WirePermission = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewCreateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Create]), new(*commandHandler.CreateHandler)),

	queryHandler.NewFindAllHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindAll, []*domain.Permission]), new(*queryHandler.FindAllHandler)),

	wire.Struct(new(fiber.CreateHandler), "*"),
	wire.Struct(new(fiber.PermissionsHandler), "*"),

	model.NewGormPermissionRepo,
	wire.Bind(new(commandHandler.CreateWriteModel), new(*model.GormPermissionRepo)),
	wire.Bind(new(queryHandler.FindAllReadModel), new(*model.GormPermissionRepo)),
)
