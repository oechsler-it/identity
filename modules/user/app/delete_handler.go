package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type DeleteWriteModel interface {
	Count(ctx context.Context) (int, error)
	Delete(ctx context.Context, id domain.UserId) error
}

type DeleteHandler struct {
	writeModel DeleteWriteModel
}

func NewDeleteHandler(writeModel DeleteWriteModel) *DeleteHandler {
	return &DeleteHandler{
		writeModel: writeModel,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd command.Delete) error {
	count, err := h.writeModel.Count(ctx)
	if err != nil {
		return err
	}
	if count == 1 {
		return domain.ErrCanNotDeleteLastUser
	}
	return h.writeModel.Delete(ctx, cmd.Id)
}
