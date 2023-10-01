package model

import "time"

type GormTokenModel struct {
	Id          string `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	OwnerUserId string
	OwnerUser   GormUserModel         `gorm:"foreignkey:OwnerUserId;constraint:OnDelete:CASCADE;"`
	Permissions []GormPermissionModel `gorm:"many2many:token_permissions;constraint:OnDelete:CASCADE;"`
	ExpiresAt   *time.Time
}

func (GormTokenModel) TableName() string {
	return "tokens"
}

type GormUserModel struct {
	Id string `gorm:"primary_key"`
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
