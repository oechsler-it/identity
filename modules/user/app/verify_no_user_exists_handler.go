package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type VerifyNoUserExistsRedModel interface {
	Count(ctx context.Context) (int, error)
}

type VerifyNoUserExistsHandler struct {
	readModel VerifyNoUserExistsRedModel
}

func NewVerifyNoUserExistsHandler(
	readModel VerifyNoUserExistsRedModel,
) *VerifyNoUserExistsHandler {
	return &VerifyNoUserExistsHandler{
		readModel: readModel,
	}
}

func (h *VerifyNoUserExistsHandler) Handle(ctx context.Context, _ command.VerifyNoUserExists) error {
	count, err := h.readModel.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return domain.ErrAUserExists
	}
	return nil
}
