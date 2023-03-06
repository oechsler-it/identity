package session

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/session/infra/model"
)
import commandHandler "github.com/oechsler-it/identity/modules/session/app"

type Options struct {
	DeviceIdMiddleware *fiber.DeviceIdMiddleware
	RenewMiddleware    *fiber.RenewMiddleware
	FiberLoginHandler  *fiber.LoginHandler
}

func UseSession(opts *Options) {
	fiber.UseDeviceIdMiddleware(opts.DeviceIdMiddleware)
	fiber.UseRenewMiddleware(opts.RenewMiddleware)
	fiber.UseFiberLoginHandler(opts.FiberLoginHandler)
}

var WireSession = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewInitiateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Initiate]), new(*commandHandler.InitiateHandler)),

	commandHandler.NewRenewHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Renew]), new(*commandHandler.RenewHandler)),

	wire.Struct(new(fiber.DeviceIdMiddleware), "*"),
	wire.Struct(new(fiber.RenewMiddleware), "*"),
	wire.Struct(new(fiber.LoginHandler), "*"),

	model.NewInMemorySessionRepo,
	wire.Bind(new(commandHandler.InitiateWriteModel), new(*model.InMemorySessionRepo)),
	wire.Bind(new(commandHandler.RenewWriteModel), new(*model.InMemorySessionRepo)),
)
