package mock

import (
	"context"
	"github.com/oechsler-it/identity/modules/user/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserModel struct {
	mock.Mock
}

func (m *MockUserModel) NextId(ctx context.Context) (domain.UserId, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.UserId), args.Error(1)
}

func (m *MockUserModel) Create(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
