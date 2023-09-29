package command

import "github.com/oechsler-it/identity/modules/permission/domain"

type VerifyPermissionNotExists struct {
	Name domain.PermissionName
}
