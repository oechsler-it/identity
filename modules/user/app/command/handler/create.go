package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type CreatePasswordService interface {
	Hash(password string) (domain.HashedPassword, error)
}

type CreateWriteModel interface {
	Create(ctx context.Context, user *domain.User) error
}

type CreateHandler struct {
	validate   *validator.Validate
	password   CreatePasswordService
	writeModel CreateWriteModel
}

func NewCreateHandler(
	validate *validator.Validate,
	password CreatePasswordService,
	writeModel CreateWriteModel,
) *CreateHandler {
	return &CreateHandler{
		validate:   validate,
		password:   password,
		writeModel: writeModel,
	}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd command.Create) error {
	hashedPassword, err := h.password.Hash(cmd.Password)
	if err != nil {
		return err
	}

	user := domain.CreateUser(
		cmd.Id,
		domain.Profile{
			FirstName: cmd.Profile.FirstName,
			LastName:  cmd.Profile.LastName,
		},
		hashedPassword,
	)

	if err := h.validate.Struct(user); err != nil {
		return err
	}

	return h.writeModel.Create(ctx, user)
}
