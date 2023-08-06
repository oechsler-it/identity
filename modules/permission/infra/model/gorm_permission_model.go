package model

import "time"

type GormPermissionModel struct {
	Name        string `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	Users       []GormUserModel `gorm:"many2many:user_permissions;"`
}

func (GormPermissionModel) TableName() string {
	return "permissions"
}

type GormUserModel struct {
	Id string `gorm:"primary_key"`
}

func (GormUserModel) TableName() string {
	return "users"
}
