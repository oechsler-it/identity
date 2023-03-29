package model

import (
	userModel "github.com/oechsler-it/identity/modules/user/infra/model"
	"time"
)

type GormSessionModel struct {
	Id            string `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	OwnerDeviceId string
	OwnerUserId   string
	OwnerUser     userModel.GormUserModel
	ExpiresAt     time.Time
	Renewable     bool
}

func (GormSessionModel) TableName() string {
	return "sessions"
}
