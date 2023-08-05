package app

import (
	"context"

	"github.com/oechsler-it/identity/modules/permission/app/query"
	"github.com/oechsler-it/identity/modules/permission/domain"
)

type FindAllReadModel interface {
	FindAll(ctx context.Context) ([]*domain.Permission, error)
}

type FindAllHandler struct {
	readModel FindAllReadModel
}

func NewFindAllHandler(readModel FindAllReadModel) *FindAllHandler {
	return &FindAllHandler{
		readModel: readModel,
	}
}

func (h *FindAllHandler) Handle(ctx context.Context, qry query.FindAll) ([]*domain.Permission, error) {
	return h.readModel.FindAll(ctx)
}
