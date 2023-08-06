package domain

import (
	"time"

	"github.com/samber/lo"
)

type User struct {
	Id             UserId         `validate:"required"`
	Profile        Profile        `validate:"required,dive"`
	HashedPassword HashedPassword `validate:"required"`
	Permissions    []Permission   `validate:"required"`
	CreatedAt      time.Time      `validate:"required"`
	UpdatedAt      time.Time      `validate:"required"`
}

// ---

func CreateUser(
	id UserId,
	profile Profile,
	hashedPassword HashedPassword,
) *User {
	return &User{
		Id:             id,
		Profile:        profile,
		HashedPassword: hashedPassword,
		Permissions:    []Permission{},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// Assertions

func (u *User) MustHavePermission(permission Permission) error {
	if !lo.Contains(u.Permissions, permission) {
		return ErrUserDoesNotHavePermission
	}
	return nil
}

// Actions

func (u *User) GrantPermission(permission Permission) {
	u.Permissions = append(u.Permissions, permission)
}

func (u *User) RemovePermission(permission Permission) {
	u.Permissions = lo.Filter(u.Permissions, func(p Permission, _ int) bool {
		return p != permission
	})
}
