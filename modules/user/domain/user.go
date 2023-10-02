package domain

import (
	"time"

	"github.com/samber/lo"
)

type User struct {
	Id             UserId         `validate:"required"`
	HashedPassword HashedPassword `validate:"required"`
	Permissions    []Permission   `validate:"required"`
	CreatedAt      time.Time      `validate:"required"`
	UpdatedAt      time.Time      `validate:"required"`
}

// Assertions

func (u *User) MustHavePermission(permission Permission) error {
	if !u.HasPermission(permission) {
		return ErrUserDoesNotHavePermission
	}
	return nil
}

func (u *User) MustHavePermissionAkinTo(permission Permission) error {
	if !u.HasPermissionAkinTo(permission) {
		return ErrUserDoesNotHavePermission
	}
	return nil
}

func (u *User) MustNotHavePermission(permission Permission) error {
	if !u.HasPermission(permission) {
		return ErrUserAlreadyHasPermission
	}
	return nil
}

// Getters

func (u *User) HasPermission(permission Permission) bool {
	return lo.Contains(u.Permissions, permission)
}

func (u *User) HasPermissionAkinTo(permission Permission) bool {
	_, found := lo.Find(u.Permissions, func(p Permission) bool {
		return p.IsAkinTo(permission)
	})
	return found
}

// Actions

func CreateUser(
	id UserId,
	hashedPassword HashedPassword,
) *User {
	return &User{
		Id:             id,
		HashedPassword: hashedPassword,
		Permissions:    []Permission{},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

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
