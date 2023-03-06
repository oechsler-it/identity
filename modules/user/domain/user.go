package domain

import "time"

type User struct {
	Id             UserId         `validate:"required"`
	Profile        Profile        `validate:"required,dive"`
	HashedPassword HashedPassword `validate:"required"`
	CreatedAt      time.Time      `validate:"required"`
	UpdatedAt      time.Time      `validate:"required"`
}

// ---

func CreateUser(
	id UserId,
	profile Profile,
	hashedPassword HashedPassword,
) *User {
	return &User{
		Id:             id,
		Profile:        profile,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
