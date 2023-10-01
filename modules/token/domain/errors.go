package domain

import "errors"

var (
	ErrTokenNotFound              = errors.New("token not found")
	ErrTokenAlreadyExists         = errors.New("token already exists")
	ErrTokenIsExpired             = errors.New("token is expired")
	ErrTokenDoesNotHavePermission = errors.New("token does not have permission")
)
