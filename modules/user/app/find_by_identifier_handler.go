package app

import (
	"context"
	"github.com/oechsler-it/identity/modules/user/app/query"
	"github.com/oechsler-it/identity/modules/user/domain"
	uuid "github.com/satori/go.uuid"
)

type FindByIdentifierReadModel interface {
	FindById(ctx context.Context, id domain.UserId) (*domain.User, error)
}

type FindByIdentifierHandler struct {
	readModel FindByIdentifierReadModel
}

func NewFindByIdentifierHandler(readModel FindByIdentifierReadModel) *FindByIdentifierHandler {
	return &FindByIdentifierHandler{
		readModel: readModel,
	}
}

func (h *FindByIdentifierHandler) Handle(ctx context.Context, query query.FindByIdentifier) (*domain.User, error) {
	// TODO: Check other valid identifiers like email, username, etc.

	identifierAsUuid, _ := uuid.FromString(query.Identifier)
	userId := domain.UserId(identifierAsUuid)
	return h.readModel.FindById(ctx, userId)
}
