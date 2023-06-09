package model

import "time"

type GormPermissionModel struct {
	Name        string `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
}

func (GormPermissionModel) TableName() string {
	return "permissions"
}
