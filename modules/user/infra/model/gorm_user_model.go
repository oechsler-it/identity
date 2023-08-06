package model

import (
	"time"
)

type GormUserModel struct {
	Id             string `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FirstName      string
	LastName       string
	HashedPassword string
	Sessions       []GormSessionModel `gorm:"foreignkey:OwnerUserId"`
}

func (GormUserModel) TableName() string {
	return "users"
}

type GormSessionModel struct {
	Id            string `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	OwnerDeviceId string
	OwnerUserId   string
	ExpiresAt     time.Time
	Renewable     bool
}

func (GormSessionModel) TableName() string {
	return "sessions"
}
