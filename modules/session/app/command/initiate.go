package command

import (
	"github.com/oechsler-it/identity/modules/session/domain"
)

type Initiate struct {
	Id                domain.SessionId
	UserId            domain.UserId
	DeviceId          domain.DeviceId
	LifetimeInSeconds int
	Renewable         bool
}
