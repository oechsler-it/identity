package mock

import (
	"github.com/oechsler-it/identity/modules/user/domain"
	"github.com/stretchr/testify/mock"
)

type MockPasswordService struct {
	mock.Mock
}

func (s *MockPasswordService) Hash(password string) (domain.HashedPassword, error) {
	args := s.Called(password)
	return args.Get(0).(domain.HashedPassword), args.Error(1)
}

func (s *MockPasswordService) Match(hashedPassword domain.HashedPassword, password string) (bool, error) {
	args := s.Called(hashedPassword, password)
	return args.Bool(0), args.Error(1)
}
