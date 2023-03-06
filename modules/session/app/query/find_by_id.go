package query

import "github.com/oechsler-it/identity/modules/session/domain"

type FindById struct {
	Id domain.SessionId
}
