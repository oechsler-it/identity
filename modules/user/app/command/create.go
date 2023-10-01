package command

import "github.com/oechsler-it/identity/modules/user/domain"

type Create struct {
	Id       domain.UserId
	Password domain.PlainPassword
}
