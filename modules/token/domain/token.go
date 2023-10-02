package domain

import (
	"github.com/samber/lo"
	"time"
)

type Token struct {
	Id          TokenId      `validate:"required"`
	CreatedAt   time.Time    `validate:"required"`
	UpdatedAt   time.Time    `validate:"required"`
	Description string       `validate:"required"`
	OwnedBy     Owner        `validate:"required"`
	Permissions []Permission `validate:"required"`
	ExpiresAt   *time.Time
}

// Assertions

func (t *Token) MustNotBeExpired() error {
	if !t.IsActive() {
		return ErrTokenIsExpired
	}
	return nil
}

func (t *Token) MustHavePermission(permission Permission) error {
	if !t.HasPermission(permission) {
		return ErrTokenDoesNotHavePermission
	}
	return nil
}

func (t *Token) MustHavePermissionAkinTo(permission Permission) error {
	if !t.HasPermissionAkinTo(permission) {
		return ErrTokenDoesNotHavePermission
	}
	return nil
}

// Getters

func (t *Token) HasPermission(permission Permission) bool {
	return lo.Contains(t.Permissions, permission)
}

func (t *Token) HasPermissionAkinTo(permission Permission) bool {
	_, found := lo.Find(t.Permissions, func(p Permission) bool {
		return p.IsAkinTo(permission)
	})
	return found
}

func (t *Token) IsActive() bool {
	return t.ExpiresAt == nil || time.Now().Before(*t.ExpiresAt)
}

// Actions

func IssueToken(
	id TokenId,
	description string,
	ownedBy Owner,
	permissions []Permission,
	expiresAt *time.Time,
) (*Token, error) {
	token := &Token{
		Id:          id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Description: description,
		OwnedBy:     ownedBy,
		Permissions: permissions,
		ExpiresAt:   expiresAt,
	}

	if err := token.MustNotBeExpired(); err != nil {
		return nil, err
	}

	return token, nil
}
