package domain_test

import (
	"github.com/oechsler-it/identity/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	userId := domain.UserId("123")
	profile := domain.Profile{
		FirstName: "John",
		LastName:  "Doe",
	}
	hashedPassword := domain.HashedPassword("hashedPassword")

	user := domain.CreateUser(
		userId,
		profile,
		hashedPassword,
	)

	assert.Equal(t, userId, user.GetId())
	assert.Equal(t, profile, user.GetProfile())
	assert.Equal(t, hashedPassword, user.GetHashedPassword())
}
