package domain

import "time"

type Permission struct {
	Name        PermissionName `validate:"required"`
	Description string
	CreatedAt   time.Time `validate:"required"`
	UpdatedAt   time.Time `validate:"required"`
}

// Actions

func CreatePermission(
	name PermissionName,
	description string,
) *Permission {
	return &Permission{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
