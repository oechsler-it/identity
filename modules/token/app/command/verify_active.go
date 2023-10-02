package command

import "github.com/oechsler-it/identity/modules/token/domain"

type VerifyActive struct {
	Id domain.TokenId
}
