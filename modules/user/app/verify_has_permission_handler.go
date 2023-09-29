package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type VerifyHasPermissionRedModel interface {
	FindById(ctx context.Context, id domain.UserId) (*domain.User, error)
}

type VerifyHasPermissionHandler struct {
	readModel VerifyHasPermissionRedModel
}

func NewVerifyHasPermissionHandler(
	readModel VerifyHasPermissionRedModel,
) *VerifyHasPermissionHandler {
	return &VerifyHasPermissionHandler{
		readModel: readModel,
	}
}

func (h *VerifyHasPermissionHandler) Handle(ctx context.Context, cmd command.VerifyHasPermission) error {
	user, err := h.readModel.FindById(ctx, cmd.Id)
	if err != nil {
		return err
	}
	return user.MustHavePermissionAkinTo(cmd.Permission)
}
