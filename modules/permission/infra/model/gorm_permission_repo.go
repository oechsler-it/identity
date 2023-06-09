package model

import (
	"context"
	"github.com/oechsler-it/identity/modules/permission/domain"
	"github.com/oechsler-it/identity/runtime"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
)

type GormPermissionRepo struct {
	database *gorm.DB
}

func NewGormPermissionRepo(
	database *gorm.DB,
	logger *logrus.Logger,
	hooks *runtime.Hooks,
) *GormPermissionRepo {
	hooks.OnStart(func(ctx context.Context) error {
		logger.WithFields(logrus.Fields{
			"name":  reflect.TypeOf(GormPermissionModel{}).Name(),
			"table": GormPermissionModel{}.TableName(),
		}).Info("Migrating model")
		return database.AutoMigrate(&GormPermissionModel{})
	})

	return &GormPermissionRepo{
		database: database,
	}
}

func (m *GormPermissionRepo) FindByName(_ context.Context, name string) (*domain.Permission, error) {
	var model GormPermissionModel
	if err := m.database.Where("name = ?", name).First(&model).Error; err != nil {
		return nil, domain.ErrPermissionNotFound
	}
	return m.toPermission(model)
}

func (m *GormPermissionRepo) Create(ctx context.Context, permission *domain.Permission) error {
	if _, err := m.FindByName(ctx, string(permission.Name)); err == nil {
		return domain.ErrPermissionAlreadyExists
	}

	model := m.toModel(permission)
	return m.database.Create(&model).Error
}

func (m *GormPermissionRepo) toPermission(model GormPermissionModel) (*domain.Permission, error) {
	return &domain.Permission{
		Name:        domain.PermissionName(model.Name),
		Description: model.Description,
	}, nil
}

func (m *GormPermissionRepo) toModel(permission *domain.Permission) GormPermissionModel {
	return GormPermissionModel{
		Name:        string(permission.Name),
		Description: permission.Description,
	}
}
