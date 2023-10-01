package model

import (
	"time"
)

type GormUserModel struct {
	Id             string `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	HashedPassword string
	Permissions    []GormPermissionModel `gorm:"many2many:user_permissions;constraint:OnDelete:CASCADE;"`
	Sessions       []GormSessionModel    `gorm:"foreignkey:OwnerUserId;"`
}

func (GormUserModel) TableName() string {
	return "users"
}

type GormPermissionModel struct {
	Name string `gorm:"primary_key"`
}

func (GormPermissionModel) TableName() string {
	return "permissions"
}

type GormSessionModel struct {
	Id          string `gorm:"primary_key"`
	OwnerUserId string
}

func (GormSessionModel) TableName() string {
	return "sessions"
}
