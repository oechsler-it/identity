package domain

import (
	"time"

	"github.com/samber/lo"
)

type User struct {
	Id             UserId         `validate:"required"`
	Profile        Profile        `validate:"required"`
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

func (u *User) MustHavePermissionAkinTo(permission Permission) error {
	_, found := lo.Find(u.Permissions, func(p Permission) bool {
		return p.IsAkinTo(permission)
	})
	if !found {
		return ErrUserDoesNotHavePermission
	}
	return nil
}

func (u *User) MustNotHavePermission(permission Permission) error {
	if lo.Contains(u.Permissions, permission) {
		return ErrUserAlreadyHasPermission
	}
	return nil
}

// Actions

func (u *User) GrantPermission(permission Permission) error {
	if err := u.MustNotHavePermission(permission); err != nil {
		return err
	}

	u.Permissions = append(u.Permissions, permission)

	return nil
}

func (u *User) RemovePermission(permission Permission) error {
	if err := u.MustHavePermission(permission); err != nil {
		return err
	}

	u.Permissions = lo.Filter(u.Permissions, func(p Permission, _ int) bool {
		return p != permission
	})

	return nil
}
