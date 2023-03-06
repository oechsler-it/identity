package command

import "github.com/oechsler-it/identity/modules/session/domain"

type Renew struct {
	Id                   domain.SessionId
	NewLifeTimeInSeconds int
}
