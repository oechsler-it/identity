package command

import "github.com/oechsler-it/identity/modules/token/domain"

type VerifyHasPermission struct {
	Id         domain.TokenIdPartial
	Permission domain.Permission
}
