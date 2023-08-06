package app

import (
	"context"

	"github.com/oechsler-it/identity/modules/permission/app/query"
	"github.com/oechsler-it/identity/modules/permission/domain"
)

type FindByNameReadModel interface {
	FindByName(ctx context.Context, name domain.PermissionName) (*domain.Permission, error)
}

type FindByNameHandler struct {
	readModel FindByNameReadModel
}

func NewFindByNameHandler(
	readModel FindByNameReadModel,
) *FindByNameHandler {
	return &FindByNameHandler{
		readModel: readModel,
	}
}

func (h *FindByNameHandler) Handle(ctx context.Context, qry query.FindByName) (*domain.Permission, error) {
	return h.readModel.FindByName(ctx, qry.Name)
}
