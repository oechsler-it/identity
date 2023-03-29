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
}

func (GormUserModel) TableName() string {
	return "users"
}
