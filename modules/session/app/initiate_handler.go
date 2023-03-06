package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
	"time"
)

type InitiateWriteModel interface {
	Create(ctx context.Context, session *domain.Session) error
}

type InitiateHandler struct {
	validate   *validator.Validate
	writeModel InitiateWriteModel
}

func NewInitiateHandler(
	validate *validator.Validate,
	writeModel InitiateWriteModel,
) *InitiateHandler {
	return &InitiateHandler{
		validate:   validate,
		writeModel: writeModel,
	}
}

func (h *InitiateHandler) Handle(ctx context.Context, cmd command.Initiate) error {
	session, err := domain.InitiateSession(
		cmd.Id,
		domain.Owner{
			UserId:   cmd.UserId,
			DeviceId: cmd.DeviceId,
		},
		time.Now().Add(time.Duration(cmd.LifetimeInSeconds)*time.Second),
		cmd.Renewable,
	)
	if err != nil {
		return err
	}

	if err := h.validate.Struct(session); err != nil {
		return err
	}

	return h.writeModel.Create(ctx, session)
}
