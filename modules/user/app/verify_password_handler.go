package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type VerifyPasswordService interface {
	Match(hashedPassword domain.HashedPassword, plainPassword domain.PlainPassword) (bool, error)
}

type VerifyPasswordReadModel interface {
	FindById(ctx context.Context, id domain.UserId) (*domain.User, error)
}

type VerifyPasswordHandler struct {
	passwordService VerifyPasswordService
	readModel       VerifyPasswordReadModel
}

func NewVerifyPasswordHandler(
	passwordService VerifyPasswordService,
	readModel VerifyPasswordReadModel,
) *VerifyPasswordHandler {
	return &VerifyPasswordHandler{
		passwordService: passwordService,
		readModel:       readModel,
	}
}

func (h *VerifyPasswordHandler) Handle(ctx context.Context, cmd command.VerifyPassword) error {
	user, err := h.readModel.FindById(ctx, cmd.Id)
	if err != nil {
		return err
	}

	match, err := h.passwordService.Match(user.HashedPassword, cmd.Password)
	if err != nil {
		return err
	}
	if !match {
		return domain.ErrInvalidPassword
	}

	return nil
}
