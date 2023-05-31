package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
)

type FindByOwnerUserIdReadModel interface {
	FindByOwnerUserId(ctx context.Context, userId domain.UserId) ([]*domain.Session, error)
}

type FindByOwnerUserIdHandler struct {
	readModel FindByOwnerUserIdReadModel
}

func NewFindByOwnerUserIdHandler(readModel FindByOwnerUserIdReadModel) *FindByOwnerUserIdHandler {
	return &FindByOwnerUserIdHandler{
		readModel: readModel,
	}
}

func (h *FindByOwnerUserIdHandler) Handle(ctx context.Context, query query.FindByOwnerUserId) ([]*domain.Session, error) {
	return h.readModel.FindByOwnerUserId(ctx, query.UserId)
}
