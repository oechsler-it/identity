package domain

type Owner struct {
	DeviceId DeviceId
	UserId   UserId `validate:"required"`
}
