package command

import "github.com/oechsler-it/identity/modules/session/domain"

type VerifyActive struct {
	Id       domain.SessionId
	DeviceId domain.DeviceId
}
