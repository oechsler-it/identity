package domain

import "errors"

var (
	ErrSessionNotFound        = errors.New("session not found")
	ErrSessionAlreadyExists   = errors.New("session already exists")
	ErrSessionMustBeRenewable = errors.New("session must be renewable")
	ErrSessionIsExpired       = errors.New("session is expired")
	ErrInvalidDeviceId        = errors.New("invalid device id")
)
