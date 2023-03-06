package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
)

type RevokeWriteModel interface {
	Delete(ctx context.Context, id domain.SessionId, handler func(session *domain.Session) error) error
}

type RevokeHandler struct {
	validate   *validator.Validate
	writeModel RevokeWriteModel
}

func NewRevokeHandler(
	validate *validator.Validate,
	writeModel RevokeWriteModel,
) *RevokeHandler {
	return &RevokeHandler{
		validate:   validate,
		writeModel: writeModel,
	}
}

func (h *RevokeHandler) Handle(ctx context.Context, cmd command.Revoke) error {
	return h.writeModel.Delete(ctx, cmd.Id, func(session *domain.Session) error {
		if err := session.Revoke(); err != nil {
			return err
		}

		if err := h.validate.Struct(session); err != nil {
			return err
		}

		return nil
	})
}
