package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
)

type VerifyActiveReadModel interface {
	FindById(ctx context.Context, id domain.SessionId) (*domain.Session, error)
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
	session, err := h.readModel.FindById(ctx, cmd.Id)
	if err != nil {
		return err
	}
	return session.MustNotBeExpired()
}
