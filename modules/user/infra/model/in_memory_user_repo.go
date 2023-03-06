package model

import (
	"context"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type InMemoryUserRepo struct {
	users map[string]*UserModel
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*UserModel),
	}
}

func (m *InMemoryUserRepo) NextId(_ context.Context) (domain.UserId, error) {
	return domain.UserId(uuid.NewV4()), nil
}

func (m *InMemoryUserRepo) FindById(_ context.Context, id domain.UserId) (*domain.User, error) {
	userId := uuid.UUID(id).String()
	if user, ok := m.users[userId]; ok {
		return m.toUser(user)
	}
	return nil, domain.ErrUserNotFound
}

func (m *InMemoryUserRepo) Create(ctx context.Context, user *domain.User) error {
	if _, err := m.FindById(ctx, user.Id); err == nil {
		return domain.ErrUserAlreadyExists
	}

	userId := uuid.UUID(user.Id).String()
	if _, ok := m.users[userId]; ok {
		return nil
	}

	m.users[userId] = m.toUserModel(user)
	return nil
}

func (m *InMemoryUserRepo) toUserModel(user *domain.User) *UserModel {
	return &UserModel{
		Id:             uuid.UUID(user.Id).String(),
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      user.UpdatedAt.Format(time.RFC3339),
		FirstName:      user.Profile.FirstName,
		LastName:       user.Profile.LastName,
		HashedPassword: string(user.HashedPassword),
	}
}

func (m *InMemoryUserRepo) toUser(model *UserModel) (*domain.User, error) {
	id, err := uuid.FromString(model.Id)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339, model.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, model.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Id:        domain.UserId(id),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Profile: domain.Profile{
			FirstName: model.FirstName,
			LastName:  model.LastName,
		},
		HashedPassword: domain.HashedPassword(model.HashedPassword),
	}, nil
}
