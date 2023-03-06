package command

import "github.com/oechsler-it/identity/modules/session/domain"

type Revoke struct {
	Id domain.SessionId
}
