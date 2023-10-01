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
	Owner       Owner        `validate:"required"`
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
	if !lo.Contains(t.Permissions, permission) {
		return ErrTokenDoesNotHavePermission
	}
	return nil
}

func (t *Token) MustHavePermissionAkinTo(permission Permission) error {
	_, found := lo.Find(t.Permissions, func(p Permission) bool {
		return p.IsAkinTo(permission)
	})
	if !found {
		return ErrTokenDoesNotHavePermission
	}
	return nil
}

// Getters

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
		Owner:       ownedBy,
		Permissions: permissions,
		ExpiresAt:   expiresAt,
	}

	if err := token.MustNotBeExpired(); err != nil {
		return nil, err
	}

	return token, nil
}
