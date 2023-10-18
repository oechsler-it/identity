package token

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/cqrs"
	commandHandler "github.com/oechsler-it/identity/modules/token/app"
	queryHandler "github.com/oechsler-it/identity/modules/token/app"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/app/query"
	"github.com/oechsler-it/identity/modules/token/domain"
	"github.com/oechsler-it/identity/modules/token/infra/fiber"
	"github.com/oechsler-it/identity/modules/token/infra/model"
)

type Options struct {
	IssueTokenHandler    *fiber.IssueTokenHandler
	ActiveTokensHandler  *fiber.ActiveTokensHandler
	TokenByIdHandler     *fiber.TokenByIdHandler
	HasPermissionHandler *fiber.HasPermissionHandler
	RevokeTokenHandler   *fiber.RevokeTokenHandler
	TokenIdMiddleware    *fiber.TokenIdMiddleware
}

func UseToken(opts *Options) {
	fiber.UseIssueTokenHandler(opts.IssueTokenHandler)
	fiber.UseActiveTokensHandler(opts.ActiveTokensHandler)
	fiber.UseTokenByIdHandler(opts.TokenByIdHandler)
	fiber.UseHasPermissionHandler(opts.HasPermissionHandler)
	fiber.UseRevokeTokenHandler(opts.RevokeTokenHandler)
	fiber.UseTokenIdMiddleware(opts.TokenIdMiddleware)
}

var WireToken = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewIssueTokenHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Issue]), new(*commandHandler.IssueHandler)),

	commandHandler.NewVerifyActiveHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyActive]), new(*commandHandler.VerifyActiveHandler)),

	commandHandler.NewVerifyHasPermissionHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyHasPermission]), new(*commandHandler.VerifyHasPermissionHandler)),

	commandHandler.NewRevokeHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Revoke]), new(*commandHandler.RevokeHandler)),

	queryHandler.NewFindByIdHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindById, *domain.Token]), new(*queryHandler.FindByIdHandler)),

	queryHandler.NewFindByIdPartialHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindByIdPartial, *domain.Token]), new(*queryHandler.FindByIdPartialHandler)),

	queryHandler.NewFindByOwnerUserIdHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindByOwnerUserId, []*domain.Token]), new(*queryHandler.FindByOwnerUserIdHandler)),

	wire.Struct(new(fiber.IssueTokenHandler), "*"),
	wire.Struct(new(fiber.ActiveTokensHandler), "*"),
	wire.Struct(new(fiber.TokenByIdHandler), "*"),
	wire.Struct(new(fiber.RevokeTokenHandler), "*"),
	wire.Struct(new(fiber.HasPermissionHandler), "*"),
	wire.Struct(new(fiber.TokenIdMiddleware), "*"),

	model.NewGormTokenRepo,
	wire.Bind(new(commandHandler.IssueWriteModel), new(*model.GormTokenRepo)),
	wire.Bind(new(commandHandler.VerifyActiveReadModel), new(*model.GormTokenRepo)),
	wire.Bind(new(commandHandler.VerifyHasPermissionReadModel), new(*model.GormTokenRepo)),
	wire.Bind(new(commandHandler.RevokeWriteModel), new(*model.GormTokenRepo)),
	wire.Bind(new(queryHandler.FindByIdReadModel), new(*model.GormTokenRepo)),
	wire.Bind(new(queryHandler.FindByIdPartialReadModel), new(*model.GormTokenRepo)),
	wire.Bind(new(queryHandler.FindByOwnerUserIdReadModel), new(*model.GormTokenRepo)),
)
