package user

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/cqrs"
	commandHandler "github.com/oechsler-it/identity/modules/user/app"
	queryHandler "github.com/oechsler-it/identity/modules/user/app"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/app/query"
	"github.com/oechsler-it/identity/modules/user/domain"
	"github.com/oechsler-it/identity/modules/user/infra/fiber"
	"github.com/oechsler-it/identity/modules/user/infra/hook"
	"github.com/oechsler-it/identity/modules/user/infra/model"
	"github.com/oechsler-it/identity/modules/user/infra/service"
)

type Options struct {
	CreateRootUser  *hook.CreateRootUser
	GrantPermission *fiber.GrantPermissionHandler
}

func UseUser(opts *Options) {
	hook.UseCreateRootUser(opts.CreateRootUser)
	fiber.UseGrantPermissionHandler(opts.GrantPermission)
}

var WireUser = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewCreateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Create]), new(*commandHandler.CreateHandler)),

	commandHandler.NewVerifyPasswordHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyPassword]), new(*commandHandler.VerifyPasswordHandler)),

	commandHandler.NewVerifyNoUserExistsHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyNoUserExists]), new(*commandHandler.VerifyNoUserExistsHandler)),

	commandHandler.NewGrantPermissionHandler,
	wire.Bind(new(cqrs.CommandHandler[command.GrantPermission]), new(*commandHandler.GrantPermissionHandler)),

	queryHandler.NewFindByIdentifierHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindByIdentifier, *domain.User]), new(*queryHandler.FindByIdentifierHandler)),

	wire.Struct(new(hook.CreateRootUser), "*"),
	wire.Struct(new(fiber.GrantPermissionHandler), "*"),

	model.NewGormUserRepo,
	wire.Bind(new(commandHandler.CreateWriteModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.VerifyPasswordReadModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.VerifyNoUserExistsRedModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.GrantPermissionWriteModel), new(*model.GormUserRepo)),
	wire.Bind(new(queryHandler.FindByIdentifierReadModel), new(*model.GormUserRepo)),

	service.NewArgon2idPasswordService,
	wire.Bind(new(commandHandler.CreatePasswordService), new(*service.Argon2idPasswordService)),
	wire.Bind(new(commandHandler.VerifyPasswordService), new(*service.Argon2idPasswordService)),
)
