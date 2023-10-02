package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/token/app/command"
	"github.com/oechsler-it/identity/modules/token/domain"
	"github.com/samber/lo"
	"time"
)

type IssueWriteModel interface {
	Create(ctx context.Context, token *domain.Token) error
}

type IssueHandler struct {
	validator  *validator.Validate
	writeModel IssueWriteModel
}

func NewIssueTokenHandler(
	validator *validator.Validate,
	writeModel IssueWriteModel,
) *IssueHandler {
	return &IssueHandler{
		validator:  validator,
		writeModel: writeModel,
	}
}

func (h *IssueHandler) Handle(ctx context.Context, cmd command.Issue) error {
	var permissions []domain.Permission
	if len(cmd.IncludedPermissions) == 0 {
		permissions = cmd.UserPermissions
	} else {
		for _, permission := range cmd.IncludedPermissions {
			if !lo.Contains(cmd.UserPermissions, permission) {
				return domain.ErrTokenCanNotBeGrantedPermission
			}
			permissions = append(permissions, permission)
		}
	}

	var expiresAt *time.Time
	if cmd.LifetimeInSeconds > 0 {
		expiresAtValue := time.Now().Add(time.Duration(cmd.LifetimeInSeconds) * time.Second)
		expiresAt = &expiresAtValue
	}

	token, err := domain.IssueToken(
		cmd.Id,
		cmd.Description,
		domain.Owner{
			UserId: cmd.UserId,
		},
		permissions,
		expiresAt,
	)
	if err != nil {
		return err
	}

	if err := h.validator.Struct(token); err != nil {
		return err
	}

	return h.writeModel.Create(ctx, token)
}
