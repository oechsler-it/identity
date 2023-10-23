package service

import (
	"context"

	"github.com/oechsler-it/identity/cqrs"
	"github.com/oechsler-it/identity/modules/session/domain"
	userCommand "github.com/oechsler-it/identity/modules/user/app/command"
	userDomain "github.com/oechsler-it/identity/modules/user/domain"
)

type UserDomainPermissionService struct {
	verifyUserHasPermission cqrs.CommandHandler[userCommand.VerifyHasPermission]
}

func NewUserDomainPermissionService(
	verifyUserHasPermission cqrs.CommandHandler[userCommand.VerifyHasPermission],
) *UserDomainPermissionService {
	return &UserDomainPermissionService{
		verifyUserHasPermission: verifyUserHasPermission,
	}
}

func (s *UserDomainPermissionService) HasPermissionAkinTo(ctx context.Context, owner domain.Owner, permission string) bool {
	if err := s.verifyUserHasPermission.Handle(ctx, userCommand.VerifyHasPermission{
		Id:         userDomain.UserId(owner.UserId),
		Permission: userDomain.Permission(permission),
	}); err != nil {
		return false
	}
	return true
}
