package command

import "github.com/oechsler-it/identity/modules/user/domain"

type Create struct {
	Id       domain.UserId
	Profile  CreateProfile `json:"profile"`
	Password string        `json:"password"`
}

type CreateProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
