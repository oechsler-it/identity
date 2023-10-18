package query

import "github.com/oechsler-it/identity/modules/token/domain"

type FindByOwnerUserId struct {
	UserId domain.UserId
}
