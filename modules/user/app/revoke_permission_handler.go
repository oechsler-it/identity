package app

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type RevokePermissionWriteModel interface {
	Revoke(ctx context.Context, id domain.UserId, handler func(user *domain.User) error) error
}

type RevokePermissionHandler struct {
	validate   *validator.Validate
	writeModel RevokePermissionWriteModel
}

func NewRevokePermissionHandler(
	validate *validator.Validate,
	writeModel RevokePermissionWriteModel,
) *RevokePermissionHandler {
	return &RevokePermissionHandler{
		validate:   validate,
		writeModel: writeModel,
	}
}

func (h *RevokePermissionHandler) Handle(ctx context.Context, cmd command.RevokePermission) error {
	return h.writeModel.Revoke(ctx, cmd.Id, func(user *domain.User) error {
		permission := domain.Permission(cmd.Permission)

		if err := user.RemovePermission(permission); err != nil {
			return err
		}

		if err := h.validate.Struct(user); err != nil {
			return err
		}

		return nil
	})
}
