package domain

import "errors"

var (
	ErrPermissionAlreadyExists = errors.New("permission already exists")
	ErrPermissionNotFound      = errors.New("permission not found")
)
