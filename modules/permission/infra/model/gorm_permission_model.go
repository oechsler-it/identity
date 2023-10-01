package model

import "time"

type GormPermissionModel struct {
	Name        string `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	Users       []GormUserModel  `gorm:"many2many:user_permissions;constraint:OnDelete:CASCADE;"`
	Tokens      []GormTokenModel `gorm:"many2many:token_permissions;constraint:OnDelete:CASCADE;"`
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

type GormTokenModel struct {
	Id string `gorm:"primary_key"`
}

func (GormTokenModel) TableName() string {
	return "tokens"
}
