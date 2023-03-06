package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
)

type FindByIdReadModel interface {
	FindById(ctx context.Context, id domain.SessionId) (*domain.Session, error)
}

type FindByIdHandler struct {
	readModel FindByIdReadModel
}

func NewFindByIdHandler(readModel FindByIdReadModel) *FindByIdHandler {
	return &FindByIdHandler{
		readModel: readModel,
	}
}

func (h *FindByIdHandler) Handle(ctx context.Context, query query.FindById) (*domain.Session, error) {
	return h.readModel.FindById(ctx, query.Id)
}
