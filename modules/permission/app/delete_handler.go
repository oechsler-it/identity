package app

import (
	"context"

	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
)

type DeleteWriteModel interface {
	Delete(ctx context.Context, name domain.PermissionName, handler func(permission *domain.Permission) error) error
}

type DeleteHandler struct {
	writeModel DeleteWriteModel
}

func NewDeleteHandler(writeModel DeleteWriteModel) *DeleteHandler {
	return &DeleteHandler{
		writeModel: writeModel,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd command.Delete) error {
	name := domain.PermissionName(cmd.Name)

	return h.writeModel.Delete(ctx, name, func(permission *domain.Permission) error {
		return nil
	})
}
