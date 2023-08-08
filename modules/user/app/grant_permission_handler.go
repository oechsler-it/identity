package app

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/domain"

	permissionQuery "github.com/oechsler-it/identity/modules/permission/app/query"
	permissionDomain "github.com/oechsler-it/identity/modules/permission/domain"
)

type GrantPermissionWriteModel interface {
	Update(ctx context.Context, id domain.UserId, handler func(user *domain.User) error) error
}

type GrantPermissionHandler struct {
	permissionFindByName cqrs.QueryHandler[permissionQuery.FindByName, *permissionDomain.Permission]
	// ---
	validate   *validator.Validate
	writeModel GrantPermissionWriteModel
}

func NewGrantPermissionHandler(
	permissionFindByName cqrs.QueryHandler[permissionQuery.FindByName, *permissionDomain.Permission],
	// ---
	validate *validator.Validate,
	writeModel GrantPermissionWriteModel,
) *GrantPermissionHandler {
	return &GrantPermissionHandler{
		permissionFindByName: permissionFindByName,
		// ---
		validate:   validate,
		writeModel: writeModel,
	}
}

func (h *GrantPermissionHandler) Handle(ctx context.Context, cmd command.GrantPermission) error {
	permissionName := permissionDomain.PermissionName(cmd.Permission)

	permission, err := h.permissionFindByName.Handle(ctx, permissionQuery.FindByName{
		Name: permissionName,
	})
	if err != nil {
		return err
	}

	return h.writeModel.Update(ctx, cmd.Id, func(user *domain.User) error {
		permission := domain.Permission(permission.Name)

		if err := user.GrantPermission(permission); err != nil {
			return err
		}

		if err := h.validate.Struct(user); err != nil {
			return err
		}

		return nil
	})
}
