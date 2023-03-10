package service

import (
	"github.com/alexedwards/argon2id"
	"github.com/oechsler-it/identity/modules/user/domain"
)

type Argon2idPasswordService struct{}

func NewArgon2idPasswordService() *Argon2idPasswordService {
	return &Argon2idPasswordService{}
}

func (s *Argon2idPasswordService) Hash(
	password domain.PlainPassword,
) (domain.HashedPassword, error) {
	hashedPassword, err := argon2id.CreateHash(string(password), argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return domain.HashedPassword(hashedPassword), nil
}

func (s *Argon2idPasswordService) Match(
	hashedPassword domain.HashedPassword,
	password domain.PlainPassword,
) (bool, error) {
	return argon2id.ComparePasswordAndHash(string(password), string(hashedPassword))
}
