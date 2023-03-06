package session

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/cqrs"
	commandHandler "github.com/oechsler-it/identity/modules/session/app"
	queryHandler "github.com/oechsler-it/identity/modules/session/app"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	"github.com/oechsler-it/identity/modules/session/infra/fiber"
	"github.com/oechsler-it/identity/modules/session/infra/model"
)

type Options struct {
	DeviceIdMiddleware  *fiber.DeviceIdMiddleware
	SessionIdMiddleware *fiber.SessionIdMiddleware
	LoginHandler        *fiber.LoginHandler
	LogoutHandler       *fiber.LogoutHandler
	SessionHandler      *fiber.SessionHandler
}

func UseSession(opts *Options) {
	fiber.UseDeviceIdMiddleware(opts.DeviceIdMiddleware)
	fiber.UseSessionIdMiddleware(opts.SessionIdMiddleware)
	fiber.UseLoginHandler(opts.LoginHandler)
	fiber.UseLogoutHandler(opts.LogoutHandler)
	fiber.UseSessionHandler(opts.SessionHandler)
}

var WireSession = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewInitiateHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Initiate]), new(*commandHandler.InitiateHandler)),

	commandHandler.NewRenewHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Renew]), new(*commandHandler.RenewHandler)),

	commandHandler.NewRevokeHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Revoke]), new(*commandHandler.RevokeHandler)),

	commandHandler.NewVerifyActiveHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyActive]), new(*commandHandler.VerifyActiveHandler)),

	queryHandler.NewFindByIdHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindById, *domain.Session]), new(*queryHandler.FindByIdHandler)),

	wire.Struct(new(fiber.DeviceIdMiddleware), "*"),
	wire.Struct(new(fiber.SessionIdMiddleware), "*"),
	wire.Struct(new(fiber.RenewMiddleware), "*"),
	wire.Struct(new(fiber.ProtectMiddleware), "*"),
	wire.Struct(new(fiber.LoginHandler), "*"),
	wire.Struct(new(fiber.LogoutHandler), "*"),
	wire.Struct(new(fiber.SessionHandler), "*"),

	model.NewInMemorySessionRepo,
	wire.Bind(new(commandHandler.InitiateWriteModel), new(*model.InMemorySessionRepo)),
	wire.Bind(new(commandHandler.RenewWriteModel), new(*model.InMemorySessionRepo)),
	wire.Bind(new(commandHandler.RevokeWriteModel), new(*model.InMemorySessionRepo)),
	wire.Bind(new(commandHandler.VerifyActiveReadModel), new(*model.InMemorySessionRepo)),
	wire.Bind(new(queryHandler.FindByIdReadModel), new(*model.InMemorySessionRepo)),
)
