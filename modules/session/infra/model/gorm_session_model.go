package model

import (
	"time"
)

type GormSessionModel struct {
	Id            string `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	OwnerDeviceId string
	OwnerUserId   string
	OwnerUser     GormUserModel `gorm:"foreignkey:OwnerUserId;constraint:OnDelete:CASCADE;"`
	ExpiresAt     time.Time
	Renewable     bool
}

func (GormSessionModel) TableName() string {
	return "sessions"
}

type GormUserModel struct {
	Id string `gorm:"primary_key"`
}

func (GormUserModel) TableName() string {
	return "users"
}
