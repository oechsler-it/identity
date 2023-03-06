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
	FiberLoginHandler *fiber.FiberLoginHandler
}

func UseSession(opts *Options) {
	fiber.UseFiberLoginHandler(opts.FiberLoginHandler)
}

var WireSession = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewInitiateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Initiate]), new(*commandHandler.InitiateHandler)),

	wire.Struct(new(fiber.FiberLoginHandler), "*"),

	model.NewInMemorySessionRepo,
	wire.Bind(new(commandHandler.InitiateWriteModel), new(*model.InMemorySessionRepo)),
)
