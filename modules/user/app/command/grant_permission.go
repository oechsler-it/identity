package command

import "github.com/oechsler-it/identity/modules/user/domain"

type GrantPermission struct {
	Id         domain.UserId
	Permission domain.Permission
}
