package command

import "github.com/oechsler-it/identity/modules/user/domain"

type VerifyPassword struct {
	Id       domain.UserId
	Password domain.PlainPassword
}
