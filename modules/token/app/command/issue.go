package command

import (
	"github.com/oechsler-it/identity/modules/token/domain"
)

type Issue struct {
	Id                  domain.TokenId
	Description         string
	UserId              domain.UserId
	UserPermissions     []domain.Permission
	IncludedPermissions []domain.Permission
	LifetimeInSeconds   int
}
