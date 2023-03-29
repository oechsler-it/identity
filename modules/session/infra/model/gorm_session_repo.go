package model

import (
	"context"
	"github.com/oechsler-it/identity/modules/session/domain"
	"github.com/oechsler-it/identity/runtime"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
)

type GormSessionRepo struct {
	database *gorm.DB
}

func NewGormSessionRepo(
	database *gorm.DB,
	logger *logrus.Logger,
	hooks *runtime.Hooks,
) *GormSessionRepo {
	hooks.OnStart(func(ctx context.Context) error {
		logger.WithFields(logrus.Fields{
			"name":  reflect.TypeOf(GormSessionModel{}).Name(),
			"table": GormSessionModel{}.TableName(),
		}).Info("Migrating model")
		return database.AutoMigrate(&GormSessionModel{})
	})

	return &GormSessionRepo{
		database: database,
	}
}

func (m *GormSessionRepo) NextId(_ context.Context) (domain.SessionId, error) {
	return domain.SessionId(uuid.NewV4()), nil
}

func (m *GormSessionRepo) FindById(_ context.Context, id domain.SessionId) (*domain.Session, error) {
	sessionId := uuid.UUID(id).String()
	var model GormSessionModel
	if err := m.database.Where("id = ?", sessionId).First(&model).Error; err != nil {
		return nil, domain.ErrSessionNotFound
	}
	return m.toSession(model)
}

func (m *GormSessionRepo) Create(ctx context.Context, session *domain.Session) error {
	if _, err := m.FindById(ctx, session.Id); err == nil {
		return domain.ErrSessionAlreadyExists
	}

	model := m.toModel(session)
	return m.database.Create(&model).Error
}

func (m *GormSessionRepo) Update(ctx context.Context, id domain.SessionId, handler func(session *domain.Session) error) error {
	session, err := m.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := handler(session); err != nil {
		return err
	}

	session.Id = id
	model := m.toModel(session)
	return m.database.Save(&model).Error
}

func (m *GormSessionRepo) Delete(ctx context.Context, id domain.SessionId, handler func(session *domain.Session) error) error {
	session, err := m.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := handler(session); err != nil {
		return err
	}

	session.Id = id
	model := m.toModel(session)
	return m.database.Delete(&model).Error
}

func (m *GormSessionRepo) toSession(model GormSessionModel) (*domain.Session, error) {
	id, err := uuid.FromString(model.Id)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.FromString(model.OwnerUserId)
	if err != nil {
		return nil, err
	}
	deviceId, err := uuid.FromString(model.OwnerDeviceId)
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		Id:        domain.SessionId(id),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		OwnedBy: domain.Owner{
			DeviceId: domain.DeviceId(deviceId),
			UserId:   domain.UserId(userId),
		},
		ExpiresAt: model.ExpiresAt,
		Renewable: model.Renewable,
	}, nil
}

func (m *GormSessionRepo) toModel(session *domain.Session) GormSessionModel {
	return GormSessionModel{
		Id:            uuid.UUID(session.Id).String(),
		CreatedAt:     session.CreatedAt,
		UpdatedAt:     session.UpdatedAt,
		OwnerDeviceId: uuid.UUID(session.OwnedBy.DeviceId).String(),
		OwnerUserId:   uuid.UUID(session.OwnedBy.UserId).String(),
		ExpiresAt:     session.ExpiresAt,
		Renewable:     session.Renewable,
	}
}
