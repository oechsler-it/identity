package model

type SessionModel struct {
	Id        string
	CreatedAt string
	UpdatedAt string
	Owner     Owner
	ExpiresAt string
	Renewable bool
}

type Owner struct {
	DeviceId string
	UserId   string
}
