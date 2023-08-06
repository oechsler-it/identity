package model

import (
	"context"
	"errors"
	"reflect"

	"github.com/oechsler-it/identity/modules/user/domain"
	"github.com/oechsler-it/identity/runtime"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormUserRepo struct {
	database *gorm.DB
}

func NewGormUserRepo(
	database *gorm.DB,
	logger *logrus.Logger,
	hooks *runtime.Hooks,
) *GormUserRepo {
	hooks.OnStart(func(ctx context.Context) error {
		logger.WithFields(logrus.Fields{
			"name":  reflect.TypeOf(GormUserModel{}).Name(),
			"table": GormUserModel{}.TableName(),
		}).Info("Migrating model")
		return database.AutoMigrate(&GormUserModel{})
	})

	return &GormUserRepo{
		database: database,
	}
}

func (m *GormUserRepo) NextId(_ context.Context) (domain.UserId, error) {
	return domain.UserId(uuid.NewV4()), nil
}

func (m *GormUserRepo) Count(ctx context.Context) (int, error) {
	var count int64
	if err := m.database.Model(&GormUserModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (m *GormUserRepo) FindById(_ context.Context, id domain.UserId) (*domain.User, error) {
	userId := uuid.UUID(id).String()
	var model GormUserModel
	if err := m.database.Where("id = ?", userId).
		Preload(clause.Associations).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return m.toUser(model)
}

func (m *GormUserRepo) Create(ctx context.Context, user *domain.User) error {
	if _, err := m.FindById(ctx, user.Id); err == nil {
		return domain.ErrUserAlreadyExists
	}

	return m.database.Create(m.toModel(user)).Error
}

func (m *GormUserRepo) toModel(user *domain.User) GormUserModel {
	return GormUserModel{
		Id:             uuid.UUID(user.Id).String(),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		FirstName:      user.Profile.FirstName,
		LastName:       user.Profile.LastName,
		HashedPassword: string(user.HashedPassword),
	}
}

func (m *GormUserRepo) toUser(model GormUserModel) (*domain.User, error) {
	id, err := uuid.FromString(model.Id)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Id:        domain.UserId(id),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Profile: domain.Profile{
			FirstName: model.FirstName,
			LastName:  model.LastName,
		},
		HashedPassword: domain.HashedPassword(model.HashedPassword),
	}, nil
}
