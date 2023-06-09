package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/permission/app/command"
	"github.com/oechsler-it/identity/modules/permission/domain"
)

type CreateWriteModel interface {
	Create(ctx context.Context, permission *domain.Permission) error
}

type CreateHandler struct {
	validate   *validator.Validate
	writeModel CreateWriteModel
}

func NewCreateHandler(
	validate *validator.Validate,
	writeModel CreateWriteModel,
) *CreateHandler {
	return &CreateHandler{
		validate:   validate,
		writeModel: writeModel,
	}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd command.Create) error {
	permission := domain.CreatePermission(
		cmd.Name,
		cmd.Description,
	)

	if err := h.validate.Struct(permission); err != nil {
		return err
	}

	return h.writeModel.Create(ctx, permission)
}
