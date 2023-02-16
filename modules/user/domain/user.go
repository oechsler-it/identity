package domain

import (
	"errors"
	"github.com/oechsler-it/identity/ddd"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserId string

type Profile struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

type HashedPassword string

type User struct {
	ddd.Entity[UserId] `validate:"required,dive"`
	Profile            Profile        `validate:"required,dive"`
	HashedPassword     HashedPassword `validate:"required"`
}

func (u *User) GetProfile() Profile {
	return u.Profile
}

func (u *User) GetHashedPassword() HashedPassword {
	return u.HashedPassword
}

func CreateUser(
	id UserId,
	profile Profile,
	hashedPassword HashedPassword,
) *User {
	return ddd.Create(id, func(e ddd.Entity[UserId]) *User {
		return &User{
			Entity:         e,
			Profile:        profile,
			HashedPassword: hashedPassword,
		}
	})
}
