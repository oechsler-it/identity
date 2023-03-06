package model

import (
	"context"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type InMemorySessionRepo struct {
	sessions map[string]*SessionModel
}

func NewInMemorySessionRepo() *InMemorySessionRepo {
	return &InMemorySessionRepo{
		sessions: make(map[string]*SessionModel),
	}
}

func (m *InMemorySessionRepo) NextId(_ context.Context) (domain.SessionId, error) {
	return domain.SessionId(uuid.NewV4()), nil
}

func (m *InMemorySessionRepo) FindById(_ context.Context, id domain.SessionId) (*domain.Session, error) {
	sessionId := uuid.UUID(id).String()
	if sessionModel, ok := m.sessions[sessionId]; ok {
		return m.toSession(sessionModel)
	}
	return nil, domain.ErrSessionNotFound
}

func (m *InMemorySessionRepo) Create(ctx context.Context, session *domain.Session) error {
	if _, err := m.FindById(ctx, session.Id); err == nil {
		return domain.ErrSessionAlreadyExists
	}

	sessionId := uuid.UUID(session.Id).String()
	m.sessions[sessionId] = m.toSessionModel(session)
	return nil
}

func (m *InMemorySessionRepo) Update(ctx context.Context, id domain.SessionId, handler func(session *domain.Session) error) error {
	session, err := m.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := handler(session); err != nil {
		return err
	}

	session.Id = id
	sessionId := uuid.UUID(id).String()
	m.sessions[sessionId] = m.toSessionModel(session)
	return nil
}

func (m *InMemorySessionRepo) Delete(ctx context.Context, id domain.SessionId, handler func(session *domain.Session) error) error {
	session, err := m.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := handler(session); err != nil {
		return err
	}

	session.Id = id
	sessionId := uuid.UUID(id).String()
	delete(m.sessions, sessionId)
	return nil
}

func (m *InMemorySessionRepo) toSessionModel(session *domain.Session) *SessionModel {
	return &SessionModel{
		Id:        uuid.UUID(session.Id).String(),
		CreatedAt: session.CreatedAt.Format(time.RFC3339),
		UpdatedAt: session.UpdatedAt.Format(time.RFC3339),
		Owner: Owner{
			DeviceId: uuid.UUID(session.OwnedBy.DeviceId).String(),
			UserId:   uuid.UUID(session.OwnedBy.UserId).String(),
		},
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
		Renewable: session.Renewable,
	}
}

func (m *InMemorySessionRepo) toSession(model *SessionModel) (*domain.Session, error) {
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
	expiresAt, err := time.Parse(time.RFC3339, model.ExpiresAt)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.FromString(model.Owner.UserId)
	if err != nil {
		return nil, err
	}
	deviceId, err := uuid.FromString(model.Owner.DeviceId)
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		Id:        domain.SessionId(id),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		OwnedBy: domain.Owner{
			UserId:   domain.UserId(userId),
			DeviceId: domain.DeviceId(deviceId),
		},
		ExpiresAt: expiresAt,
		Renewable: model.Renewable,
	}, nil
}
