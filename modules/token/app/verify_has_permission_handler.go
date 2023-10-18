package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
)

type VerifyHasPermissionReadModel interface {
	FindByIdPartial(ctx context.Context, id domain.TokenIdPartial) (*domain.Token, error)
}

type VerifyHasPermissionHandler struct {
	readModel VerifyHasPermissionReadModel
}

func NewVerifyHasPermissionHandler(
	readModel VerifyHasPermissionReadModel,
) *VerifyHasPermissionHandler {
	return &VerifyHasPermissionHandler{
		readModel: readModel,
	}
}

func (h *VerifyHasPermissionHandler) Handle(ctx context.Context, cmd command.VerifyHasPermission) error {
	token, err := h.readModel.FindByIdPartial(ctx, cmd.Id)
	if err != nil {
		return err
	}
	if err := token.MustNotBeExpired(); err != nil {
		return err
	}
	return token.MustHavePermissionAkinTo(cmd.Permission)
}
