package model

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/oechsler-it/identity/modules/token/domain"
	"github.com/oechsler-it/identity/runtime"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"reflect"
)

type GormTokenRepo struct {
	database *gorm.DB
}

func NewGormTokenRepo(
	database *gorm.DB,
	logger *logrus.Logger,
	hooks *runtime.Hooks,
) *GormTokenRepo {
	hooks.OnStart(func(ctx context.Context) error {
		logger.WithFields(logrus.Fields{
			"name":  reflect.TypeOf(GormTokenModel{}).Name(),
			"table": GormTokenModel{}.TableName(),
		}).Info("Migrating model")
		return database.AutoMigrate(&GormTokenModel{})
	})

	return &GormTokenRepo{
		database: database,
	}
}

func (m *GormTokenRepo) NextId(_ context.Context) (domain.TokenId, error) {
	value := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, value); err != nil {
		return "", err
	}
	encodedValue := hex.EncodeToString(value)
	return domain.TokenId(encodedValue), nil
}

func (m *GormTokenRepo) FindById(_ context.Context, id domain.TokenId) (*domain.Token, error) {
	tokenId := string(id)
	var model GormTokenModel
	if err := m.database.Where("id = ?", tokenId).
		Preload(clause.Associations).
		First(&model).Error; err != nil {
		return nil, domain.ErrTokenNotFound
	}
	return m.toToken(model)
}

func (m *GormTokenRepo) FindByOwnerUserId(_ context.Context, userId domain.UserId) ([]*domain.Token, error) {
	userIdString := uuid.UUID(userId).String()
	var models []GormTokenModel
	if err := m.database.Where("owner_user_id = ?", userIdString).
		Preload(clause.Associations).
		Find(&models).Error; err != nil {
		return nil, err
	}
	tokens := make([]*domain.Token, len(models))
	for i, model := range models {
		token, err := m.toToken(model)
		if err != nil {
			return nil, err
		}
		tokens[i] = token
	}
	return tokens, nil
}

func (m *GormTokenRepo) Create(ctx context.Context, token *domain.Token) error {
	if _, err := m.FindById(ctx, token.Id); err == nil {
		return domain.ErrTokenAlreadyExists
	}

	model := m.toModel(token)
	return m.database.Create(&model).Error
}

func (m *GormTokenRepo) Update(ctx context.Context, id domain.TokenId, handler func(token *domain.Token) error) error {
	token, err := m.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := handler(token); err != nil {
		return err
	}

	token.Id = id
	model := m.toModel(token)
	return m.database.Save(&model).Error
}

func (m *GormTokenRepo) Delete(ctx context.Context, id domain.TokenId, handler func(token *domain.Token) error) error {
	token, err := m.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := handler(token); err != nil {
		return err
	}

	token.Id = id
	model := m.toModel(token)
	return m.database.Delete(&model).Error
}

func (m *GormTokenRepo) toToken(model GormTokenModel) (*domain.Token, error) {
	userId, err := uuid.FromString(model.OwnerUserId)
	if err != nil {
		return nil, err
	}

	return &domain.Token{
		Id:          domain.TokenId(model.Id),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		Description: model.Description,
		OwnedBy: domain.Owner{
			UserId: domain.UserId(userId),
		},
		Permissions: lo.Map(model.Permissions, func(model GormPermissionModel, _ int) domain.Permission {
			return domain.Permission(model.Name)
		}),
		ExpiresAt: model.ExpiresAt,
	}, nil
}

func (m *GormTokenRepo) toModel(token *domain.Token) GormTokenModel {
	return GormTokenModel{
		Id:          string(token.Id),
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
		Description: token.Description,
		OwnerUserId: uuid.UUID(token.OwnedBy.UserId).String(),
		Permissions: lo.Map(token.Permissions, func(permission domain.Permission, _ int) GormPermissionModel {
			return GormPermissionModel{
				Name: string(permission),
			}
		}),
		ExpiresAt: token.ExpiresAt,
	}
}
