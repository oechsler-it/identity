package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/token/app/query"
	"github.com/oechsler-it/identity/modules/token/domain"
)

type FindByIdPartialReadModel interface {
	FindByIdPartial(ctx context.Context, partial domain.TokenIdPartial) (*domain.Token, error)
}

type FindByIdPartialHandler struct {
	readModel FindByIdPartialReadModel
}

func NewFindByIdPartialHandler(
	readModel FindByIdPartialReadModel,
) *FindByIdPartialHandler {
	return &FindByIdPartialHandler{
		readModel: readModel,
	}
}

func (h *FindByIdPartialHandler) Handle(ctx context.Context, qry query.FindByIdPartial) (*domain.Token, error) {
	return h.readModel.FindByIdPartial(ctx, qry.IdPartial)
}
