package query

import "github.com/oechsler-it/identity/modules/permission/domain"

type FindByName struct {
	Name domain.PermissionName
}
