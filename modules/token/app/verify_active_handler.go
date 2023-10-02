package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
)

type VerifyActiveReadModel interface {
	FindById(ctx context.Context, id domain.TokenId) (*domain.Token, error)
}

type VerifyActiveHandler struct {
	readModel VerifyActiveReadModel
}

func NewVerifyActiveHandler(
	readModel VerifyActiveReadModel,
) *VerifyActiveHandler {
	return &VerifyActiveHandler{
		readModel: readModel,
	}
}

func (h *VerifyActiveHandler) Handle(ctx context.Context, cmd command.VerifyActive) error {
	token, err := h.readModel.FindById(ctx, cmd.Id)
	if err != nil {
		return err
	}
	return token.MustNotBeExpired()
}
