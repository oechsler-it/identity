package app

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/session/app/command"
	"github.com/oechsler-it/identity/modules/session/domain"
)

type UserPermissionService interface {
	HasPermissionAkinTo(ctx context.Context, owner domain.Owner, permission string) bool
}

type RevokeWriteModel interface {
	Delete(ctx context.Context, id domain.SessionId, handler func(session *domain.Session) error) error
}

type RevokeHandler struct {
	validate              *validator.Validate
	writeModel            RevokeWriteModel
	userPermissionService UserPermissionService
}

func NewRevokeHandler(
	validate *validator.Validate,
	writeModel RevokeWriteModel,
	userPermissionService UserPermissionService,
) *RevokeHandler {
	return &RevokeHandler{
		validate:              validate,
		writeModel:            writeModel,
		userPermissionService: userPermissionService,
	}
}

func (h *RevokeHandler) Handle(ctx context.Context, cmd command.Revoke) error {
	return h.writeModel.Delete(ctx, cmd.Id, func(session *domain.Session) error {
		allowedByPermission := h.userPermissionService.HasPermissionAkinTo(ctx, session.OwnedBy, "all:session:revoke")

		if err := session.Revoke(cmd.RevokingEntity, allowedByPermission); err != nil {
			return err
		}

		if err := h.validate.Struct(session); err != nil {
			return err
		}

		return nil
	})
}
