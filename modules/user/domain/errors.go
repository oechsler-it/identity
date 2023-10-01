package domain

import "errors"

var (
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrAUserExists               = errors.New("a user exists")
	ErrUserDoesNotHavePermission = errors.New("user does not have permission")
	ErrUserAlreadyHasPermission  = errors.New("user already has permission")
	ErrCanNotDeleteLastUser      = errors.New("can not delete last user")
)
