package domain

import "strings"

type Permission string

var PermissionNone = Permission("")

// IsAkinTo A permission is akin to another if it is a subset of the other.
// For example, a user with the permission "user:read" is akin to
// a user with the permission "user:read:own".
func (p Permission) IsAkinTo(other Permission) bool {
	return p == other || strings.HasPrefix(string(other), string(p))
}
