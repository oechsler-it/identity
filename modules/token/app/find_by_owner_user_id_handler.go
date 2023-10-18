package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/token/app/query"
	"github.com/oechsler-it/identity/modules/token/domain"
	"github.com/samber/lo"
)

type FindByOwnerUserIdReadModel interface {
	FindByOwnerUserId(ctx context.Context, userId domain.UserId) ([]*domain.Token, error)
}

type FindByOwnerUserIdHandler struct {
	readModel FindByOwnerUserIdReadModel
}

func NewFindByOwnerUserIdHandler(readModel FindByOwnerUserIdReadModel) *FindByOwnerUserIdHandler {
	return &FindByOwnerUserIdHandler{
		readModel: readModel,
	}
}

func (h *FindByOwnerUserIdHandler) Handle(ctx context.Context, qry query.FindByOwnerUserId) ([]*domain.Token, error) {
	tokens, err := h.readModel.FindByOwnerUserId(ctx, qry.UserId)
	if err != nil {
		return nil, err
	}
	return lo.Filter(tokens, func(token *domain.Token, _ int) bool {
		return token.IsActive()
	}), nil
}
