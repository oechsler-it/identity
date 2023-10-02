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
	DeviceIdMiddleware    *fiber.DeviceIdMiddleware
	SessionIdMiddleware   *fiber.SessionMiddleware
	LoginHandler          *fiber.LoginHandler
	LogoutHandler         *fiber.LogoutHandler
	RevokeSessionHandler  *fiber.RevokeSessionHandler
	ActiveSessionsHandler *fiber.ActiveSessionsHandler
	ActiveSessionHandler  *fiber.ActiveSessionHandler
	SessionByIdHandler    *fiber.SessionByIdHandler
}

func UseSession(opts *Options) {
	fiber.UseDeviceIdMiddleware(opts.DeviceIdMiddleware)
	fiber.UseSessionMiddleware(opts.SessionIdMiddleware)
	fiber.UseLoginHandler(opts.LoginHandler)
	fiber.UseLogoutHandler(opts.LogoutHandler)
	fiber.UseRevokeSessionHandler(opts.RevokeSessionHandler)
	fiber.UseActiveSessionsHandler(opts.ActiveSessionsHandler)
	fiber.UseActiveSessionHandler(opts.ActiveSessionHandler)
	fiber.UseSessionByIdHandler(opts.SessionByIdHandler)
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

	queryHandler.NewFindByOwnerUserIdHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindByOwnerUserId, []*domain.Session]), new(*queryHandler.FindByOwnerUserIdHandler)),

	queryHandler.NewFindByIdHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindById, *domain.Session]), new(*queryHandler.FindByIdHandler)),

	wire.Struct(new(fiber.DeviceIdMiddleware), "*"),
	wire.Struct(new(fiber.SessionMiddleware), "*"),
	wire.Struct(new(fiber.RenewMiddleware), "*"),
	wire.Struct(new(fiber.ProtectSessionMiddleware), "*"),
	wire.Struct(new(fiber.LoginHandler), "*"),
	wire.Struct(new(fiber.LogoutHandler), "*"),
	wire.Struct(new(fiber.RevokeSessionHandler), "*"),
	wire.Struct(new(fiber.ActiveSessionsHandler), "*"),
	wire.Struct(new(fiber.ActiveSessionHandler), "*"),
	wire.Struct(new(fiber.SessionByIdHandler), "*"),

	model.NewGormSessionRepo,
	wire.Bind(new(commandHandler.InitiateWriteModel), new(*model.GormSessionRepo)),
	wire.Bind(new(commandHandler.RenewWriteModel), new(*model.GormSessionRepo)),
	wire.Bind(new(commandHandler.RevokeWriteModel), new(*model.GormSessionRepo)),
	wire.Bind(new(commandHandler.VerifyActiveReadModel), new(*model.GormSessionRepo)),
	wire.Bind(new(queryHandler.FindByOwnerUserIdReadModel), new(*model.GormSessionRepo)),
	wire.Bind(new(queryHandler.FindByIdReadModel), new(*model.GormSessionRepo)),
)
