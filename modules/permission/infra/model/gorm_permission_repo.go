package model

import (
	"context"
	"reflect"

	"github.com/oechsler-it/identity/modules/permission/domain"
	"github.com/oechsler-it/identity/runtime"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func (m *GormPermissionRepo) FindAll(_ context.Context) ([]*domain.Permission, error) {
	var models []GormPermissionModel
	if err := m.database.Find(&models).Error; err != nil {
		return nil, err
	}
	permissions := make([]*domain.Permission, len(models))
	for i, model := range models {
		permission, err := m.toPermission(model)
		if err != nil {
			return nil, err
		}
		permissions[i] = permission
	}
	return permissions, nil
}

func (m *GormPermissionRepo) FindByName(_ context.Context, name domain.PermissionName) (*domain.Permission, error) {
	var model GormPermissionModel
	if err := m.database.Where("name = ?", name).First(&model).Error; err != nil {
		return nil, domain.ErrPermissionNotFound
	}
	return m.toPermission(model)
}

func (m *GormPermissionRepo) Create(ctx context.Context, permission *domain.Permission) error {
	if _, err := m.FindByName(ctx, permission.Name); err == nil {
		return domain.ErrPermissionAlreadyExists
	}

	model := m.toModel(permission)
	return m.database.Create(&model).Error
}

func (m *GormPermissionRepo) Update(ctx context.Context, name domain.PermissionName, handler func(permission *domain.Permission) error) error {
	permission, err := m.FindByName(ctx, name)
	if err != nil {
		return err
	}

	if err := handler(permission); err != nil {
		return err
	}

	permission.Name = name
	model := m.toModel(permission)
	return m.database.Save(&model).Error
}

func (m *GormPermissionRepo) Delete(ctx context.Context, name domain.PermissionName, handler func(permission *domain.Permission) error) error {
	permission, err := m.FindByName(ctx, name)
	if err != nil {
		return err
	}

	if err := handler(permission); err != nil {
		return err
	}

	permission.Name = name
	model := m.toModel(permission)
	return m.database.Delete(&model).Error
}

func (m *GormPermissionRepo) toPermission(model GormPermissionModel) (*domain.Permission, error) {
	return &domain.Permission{
		Name:        domain.PermissionName(model.Name),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		Description: model.Description,
	}, nil
}

func (m *GormPermissionRepo) toModel(permission *domain.Permission) GormPermissionModel {
	return GormPermissionModel{
		Name:        string(permission.Name),
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
		Description: permission.Description,
	}
}
