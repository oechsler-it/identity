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
	fiberMiddleware "github.com/oechsler-it/identity/modules/user/infra/fiber/middleware"
	"github.com/oechsler-it/identity/modules/user/infra/hook"
	"github.com/oechsler-it/identity/modules/user/infra/model"
	"github.com/oechsler-it/identity/modules/user/infra/service"
)

type Options struct {
	CreateRootUser   *hook.CreateRootUser
	CreateUser       *fiber.CreateUserHandler
	DeleteMe         *fiber.DeleteMeHandler
	DeleteUser       *fiber.DeleteUserHandler
	Me               *fiber.MeHandler
	UserById         *fiber.UserByIdHandler
	GrantPermission  *fiber.GrantPermissionHandler
	RevokePermission *fiber.RevokePermissionHandler
	HasPermission    *fiber.HasPermissionHandler
}

func UseUser(opts *Options) {
	hook.UseCreateRootUser(opts.CreateRootUser)
	fiber.UseCreateUserHandler(opts.CreateUser)
	fiber.UseDeleteMeHandler(opts.DeleteMe)
	fiber.UseDeleteUserHandler(opts.DeleteUser)
	fiber.UseMeHandler(opts.Me)
	fiber.UseUserByIdHandler(opts.UserById)
	fiber.UseGrantPermissionHandler(opts.GrantPermission)
	fiber.UseRevokePermissionHandler(opts.RevokePermission)
	fiber.UseHasPermissionHandler(opts.HasPermission)
}

var WireUser = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewCreateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Create]), new(*commandHandler.CreateHandler)),

	commandHandler.NewDeleteHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Delete]), new(*commandHandler.DeleteHandler)),

	commandHandler.NewVerifyPasswordHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyPassword]), new(*commandHandler.VerifyPasswordHandler)),

	commandHandler.NewVerifyNoUserExistsHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyNoUserExists]), new(*commandHandler.VerifyNoUserExistsHandler)),

	commandHandler.NewGrantPermissionHandler,
	wire.Bind(new(cqrs.CommandHandler[command.GrantPermission]), new(*commandHandler.GrantPermissionHandler)),

	commandHandler.NewVerifyHasPermissionHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyHasPermission]), new(*commandHandler.VerifyHasPermissionHandler)),

	commandHandler.NewRevokePermissionHandler,
	wire.Bind(new(cqrs.CommandHandler[command.RevokePermission]), new(*commandHandler.RevokePermissionHandler)),

	queryHandler.NewFindByIdentifierHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindByIdentifier, *domain.User]), new(*queryHandler.FindByIdentifierHandler)),

	wire.Struct(new(hook.CreateRootUser), "*"),
	wire.Struct(new(fiber.CreateUserHandler), "*"),
	wire.Struct(new(fiber.DeleteUserHandler), "*"),
	wire.Struct(new(fiber.DeleteMeHandler), "*"),
	wire.Struct(new(fiber.MeHandler), "*"),
	wire.Struct(new(fiber.UserByIdHandler), "*"),
	wire.Struct(new(fiber.GrantPermissionHandler), "*"),
	wire.Struct(new(fiber.RevokePermissionHandler), "*"),
	wire.Struct(new(fiber.HasPermissionHandler), "*"),
	wire.Struct(new(fiberMiddleware.UserMiddleware), "*"),
	wire.Struct(new(fiberMiddleware.UserPermissionMiddleware), "*"),

	model.NewGormUserRepo,
	wire.Bind(new(commandHandler.CreateWriteModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.DeleteWriteModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.VerifyPasswordReadModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.VerifyNoUserExistsRedModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.GrantPermissionWriteModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.RevokePermissionWriteModel), new(*model.GormUserRepo)),
	wire.Bind(new(commandHandler.VerifyHasPermissionRedModel), new(*model.GormUserRepo)),
	wire.Bind(new(queryHandler.FindByIdentifierReadModel), new(*model.GormUserRepo)),

	service.NewArgon2idPasswordService,
	wire.Bind(new(commandHandler.CreatePasswordService), new(*service.Argon2idPasswordService)),
	wire.Bind(new(commandHandler.VerifyPasswordService), new(*service.Argon2idPasswordService)),
)
