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
	IssueTokenHandler *fiber.IssueTokenHandler
	TokenIdMiddleware *fiber.TokenIdMiddleware
}

func UseToken(opts *Options) {
	fiber.UseIssueTokenHandler(opts.IssueTokenHandler)
	fiber.UseTokenIdMiddleware(opts.TokenIdMiddleware)
}

var WireToken = wire.NewSet(
	wire.Struct(new(Options), "*"),

	commandHandler.NewIssueTokenHandler,
	wire.Bind(new(cqrs.CommandHandler[command.Issue]), new(*commandHandler.IssueHandler)),

	commandHandler.NewVerifyActiveHandler,
	wire.Bind(new(cqrs.CommandHandler[command.VerifyActive]), new(*commandHandler.VerifyActiveHandler)),

	queryHandler.NewFindByIdHandler,
	wire.Bind(new(cqrs.QueryHandler[query.FindById, *domain.Token]), new(*queryHandler.FindByIdHandler)),

	wire.Struct(new(fiber.IssueTokenHandler), "*"),
	wire.Struct(new(fiber.TokenIdMiddleware), "*"),

	model.NewGormTokenRepo,
	wire.Bind(new(commandHandler.IssueWriteModel), new(*model.GormTokenRepo)),
	wire.Bind(new(commandHandler.VerifyActiveReadModel), new(*model.GormTokenRepo)),
	wire.Bind(new(queryHandler.FindByIdReadModel), new(*model.GormTokenRepo)),
)
