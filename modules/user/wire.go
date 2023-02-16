package user

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/user/app/command"
	commandHandler "github.com/oechsler-it/identity/modules/user/app/command/handler"
	"github.com/oechsler-it/identity/modules/user/infra/hook"
	"github.com/oechsler-it/identity/modules/user/infra/model"
	"github.com/oechsler-it/identity/modules/user/infra/service"
)

type Options struct {
	HookCreate *hook.HooksCreateRootUser
}

func UseUser(opts *Options) {
	hook.UseHooksCreateRootUser(opts.HookCreate)
}

var WireUser = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewCreateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Create]), new(*commandHandler.CreateHandler)),

	wire.Struct(new(hook.HooksCreateRootUser), "*"),

	model.NewInMemoryUserModel,
	wire.Bind(new(commandHandler.CreateWriteModel), new(*model.InMemoryUserModel)),

	service.NewArgon2idPasswordService,
	wire.Bind(new(commandHandler.CreatePasswordService), new(*service.Argon2idPasswordService)),
)
