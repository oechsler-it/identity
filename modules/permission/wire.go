package permission

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/cqrs"
	commandHandler "github.com/oechsler-it/identity/modules/permission/app"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/infra/fiber"
	"github.com/oechsler-it/identity/modules/permission/infra/model"
)

type Options struct {
	CreateHandler *fiber.CreateHandler
}

func UsePermission(opts *Options) {
	fiber.UseCreateHandler(opts.CreateHandler)
}

var WirePermission = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewCreateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Create]), new(*commandHandler.CreateHandler)),

	wire.Struct(new(fiber.CreateHandler), "*"),

	model.NewGormPermissionRepo,
	wire.Bind(new(commandHandler.CreateWriteModel), new(*model.GormPermissionRepo)),
)
