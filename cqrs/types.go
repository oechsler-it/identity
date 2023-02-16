package cqrs

import "context"

type CommandHandler[TCommand any] interface {
	Handle(ctx context.Context, cmd TCommand) error
}

type QueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}
