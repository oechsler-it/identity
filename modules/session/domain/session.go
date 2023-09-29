package domain

import (
	"time"
)

type Session struct {
	Id        SessionId `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	OwnedBy   Owner     `validate:"required"`
	ExpiresAt time.Time `validate:"required"`
	Renewable bool
}

// Assertions

func (s *Session) MustNotBeExpired() error {
	if !s.IsActive() {
		return ErrSessionIsExpired
	}
	return nil
}

func (s *Session) MustBeOwnedBy(owner Owner) error {
	if s.OwnedBy != owner {
		return ErrSessionDoesNotBelongToOwner
	}
	return nil
}

func (s *Session) MustBeRenewable() error {
	if !s.Renewable {
		return ErrSessionMustBeRenewable
	}
	return nil
}

// Getters

func (s *Session) IsActive() bool {
	return time.Now().Before(s.ExpiresAt)
}

// Actions

func InitiateSession(
	id SessionId,
	ownedBy Owner,
	expiresAt time.Time,
	renewable bool,
) (*Session, error) {
	session := &Session{
		Id:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		OwnedBy:   ownedBy,
		ExpiresAt: expiresAt,
		Renewable: renewable,
	}

	if err := session.MustNotBeExpired(); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Session) Renew(expiresAt time.Time) error {
	if err := s.MustBeRenewable(); err != nil {
		return err
	}
	if err := s.MustNotBeExpired(); err != nil {
		return err
	}

	s.ExpiresAt = expiresAt
	if err := s.MustNotBeExpired(); err != nil {
		return err
	}

	s.UpdatedAt = time.Now()

	return nil
}

func (s *Session) Revoke(revokingEntity Owner) error {
	if err := s.MustBeOwnedBy(revokingEntity); err != nil {
		return err
	}
	if err := s.MustNotBeExpired(); err != nil {
		return err
	}

	s.ExpiresAt = time.Now()
	s.UpdatedAt = time.Now()

	return nil
}
