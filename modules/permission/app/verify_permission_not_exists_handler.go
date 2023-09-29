package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
)

type VerifyPermissionNotExistsReadModel interface {
	FindByName(ctx context.Context, name domain.PermissionName) (*domain.Permission, error)
}

type VerifyPermissionNotExistsHandler struct {
	readModel VerifyPermissionNotExistsReadModel
}

func NewVerifyPermissionNotExistsHandler(
	readModel VerifyPermissionNotExistsReadModel,
) *VerifyPermissionNotExistsHandler {
	return &VerifyPermissionNotExistsHandler{
		readModel: readModel,
	}
}

func (h *VerifyPermissionNotExistsHandler) Handle(ctx context.Context, cmd command.VerifyPermissionNotExists) error {
	if _, err := h.readModel.FindByName(ctx, cmd.Name); err != nil {
		return nil
	}
	return domain.ErrPermissionAlreadyExists
}
