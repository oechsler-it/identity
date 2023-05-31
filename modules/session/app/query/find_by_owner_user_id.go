package query

import "github.com/oechsler-it/identity/modules/session/domain"

type FindByOwnerUserId struct {
	UserId domain.UserId
}
