package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
)

type RevokeWriteModel interface {
	DeleteByIdPartial(ctx context.Context, idPartial domain.TokenIdPartial, handler func(token *domain.Token) error) error
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
	return h.writeModel.DeleteByIdPartial(ctx, cmd.IdPartial, func(token *domain.Token) error {
		if err := token.Revoke(cmd.RevokingEntity); err != nil {
			return err
		}

		if err := h.validate.Struct(token); err != nil {
			return err
		}

		return nil
	})
}
