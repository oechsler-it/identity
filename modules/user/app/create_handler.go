package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type CreatePasswordService interface {
	Hash(password domain.PlainPassword) (domain.HashedPassword, error)
}

type CreateWriteModel interface {
	Create(ctx context.Context, user *domain.User) error
}

type CreateHandler struct {
	validate        *validator.Validate
	passwordService CreatePasswordService
	writeModel      CreateWriteModel
}

func NewCreateHandler(
	validate *validator.Validate,
	passwordService CreatePasswordService,
	writeModel CreateWriteModel,
) *CreateHandler {
	return &CreateHandler{
		validate:        validate,
		passwordService: passwordService,
		writeModel:      writeModel,
	}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd command.Create) error {
	hashedPassword, err := h.passwordService.Hash(cmd.Password)
	if err != nil {
		return err
	}

	user := domain.CreateUser(
		cmd.Id,
		cmd.Profile,
		hashedPassword,
	)

	if err := h.validate.Struct(user); err != nil {
		return err
	}

	return h.writeModel.Create(ctx, user)
}
