package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	"time"
)

type RenewWriteModel interface {
	Update(ctx context.Context, id domain.SessionId, handler func(session *domain.Session) error) error
}

type RenewHandler struct {
	validate   *validator.Validate
	writeModel RenewWriteModel
}

func NewRenewHandler(
	validate *validator.Validate,
	writeModel RenewWriteModel,
) *RenewHandler {
	return &RenewHandler{
		validate:   validate,
		writeModel: writeModel,
	}
}

func (h *RenewHandler) Handle(ctx context.Context, cmd command.Renew) error {
	return h.writeModel.Update(ctx, cmd.Id, func(session *domain.Session) error {
		newExpiresAt := time.Now().Add(time.Duration(cmd.NewLifeTimeInSeconds) * time.Second)

		if err := session.Renew(newExpiresAt); err != nil {
			return err
		}

		if err := h.validate.Struct(session); err != nil {
			return err
		}

		return nil
	})
}
