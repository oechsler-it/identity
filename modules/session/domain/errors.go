package domain

import "errors"

var (
	ErrSessionMustBeRenewable = errors.New("session must be renewable")
	ErrSessionIsExpired       = errors.New("session is expired")
)
