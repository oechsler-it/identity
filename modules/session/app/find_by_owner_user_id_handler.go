package app

import (
	"context"

	"github.com/oechsler-it/identity/modules/session/app/query"
	"github.com/oechsler-it/identity/modules/session/domain"
	"github.com/samber/lo"
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
	sessions, err := h.readModel.FindByOwnerUserId(ctx, query.UserId)
	if err != nil {
		return nil, err
	}
	return lo.Filter(sessions, func(session *domain.Session, _ int) bool {
		return session.IsActive()
	}), nil
}
