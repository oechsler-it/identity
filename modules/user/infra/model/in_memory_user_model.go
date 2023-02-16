package model

import (
	"context"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
)

type InMemoryUserModel struct {
	users map[string]*domain.User
}

func NewInMemoryUserModel() *InMemoryUserModel {
	return &InMemoryUserModel{
		users: make(map[string]*domain.User),
	}
}

func (m *InMemoryUserModel) NextId(_ context.Context) (domain.UserId, error) {
	return domain.UserId(uuid.NewV4().String()), nil
}

func (m *InMemoryUserModel) Create(_ context.Context, user *domain.User) error {
	if _, ok := m.users[string(user.GetId())]; ok {
		return domain.ErrUserAlreadyExists
	}
	m.users[string(user.GetId())] = user
	return nil
}
