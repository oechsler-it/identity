package command

import "github.com/oechsler-it/identity/modules/permission/domain"

type Create struct {
	Name        domain.PermissionName
	Description string
}
