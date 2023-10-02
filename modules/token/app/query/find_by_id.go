package query

import "github.com/oechsler-it/identity/modules/token/domain"

type FindById struct {
	Id domain.TokenId
}
