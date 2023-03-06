package user

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/cqrs"
	commandHandler "github.com/oechsler-it/identity/modules/user/app"
	queryHandler "github.com/oechsler-it/identity/modules/user/app"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/app/query"
	"github.com/oechsler-it/identity/modules/user/domain"
	"github.com/oechsler-it/identity/modules/user/infra/hook"
	"github.com/oechsler-it/identity/modules/user/infra/model"
	"github.com/oechsler-it/identity/modules/user/infra/service"
)

type Options struct {
	CreateRootUser *hook.CreateRootUser
}

func UseUser(opts *Options) {
	hook.UseCreateRootUser(opts.CreateRootUser)
}

var WireUser = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewCreateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Create]), new(*commandHandler.CreateHandler)),

	commandHandler.NewVerifyPasswordHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyPassword]), new(*commandHandler.VerifyPasswordHandler)),

	queryHandler.NewFindByIdentifierHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindByIdentifier, *domain.User]), new(*queryHandler.FindByIdentifierHandler)),

	wire.Struct(new(hook.CreateRootUser), "*"),

	model.NewInMemoryUserRepo,
	wire.Bind(new(commandHandler.CreateWriteModel), new(*model.InMemoryUserRepo)),
	wire.Bind(new(commandHandler.VerifyPasswordReadModel), new(*model.InMemoryUserRepo)),
	wire.Bind(new(queryHandler.FindByIdentifierReadModel), new(*model.InMemoryUserRepo)),

	service.NewArgon2idPasswordService,
	wire.Bind(new(commandHandler.CreatePasswordService), new(*service.Argon2idPasswordService)),
	wire.Bind(new(commandHandler.VerifyPasswordService), new(*service.Argon2idPasswordService)),
)
