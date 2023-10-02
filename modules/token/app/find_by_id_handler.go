package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/token/app/query"
	"github.com/oechsler-it/identity/modules/token/domain"
)

type FindByIdReadModel interface {
	FindById(ctx context.Context, id domain.TokenId) (*domain.Token, error)
}

type FindByIdHandler struct {
	readModel FindByIdReadModel
}

func NewFindByIdHandler(
	readModel FindByIdReadModel,
) *FindByIdHandler {
	return &FindByIdHandler{
		readModel: readModel,
	}
}

func (h *FindByIdHandler) Handle(ctx context.Context, qry query.FindById) (*domain.Token, error) {
	return h.readModel.FindById(ctx, qry.Id)
}
