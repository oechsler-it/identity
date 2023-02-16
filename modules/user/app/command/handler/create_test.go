package handler_test

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/modules/user/app/command"
	"github.com/oechsler-it/identity/modules/user/app/command/handler"
	"github.com/oechsler-it/identity/modules/user/domain"
	modelMock "github.com/oechsler-it/identity/modules/user/infra/model/mock"
	serviceMock "github.com/oechsler-it/identity/modules/user/infra/service/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateHandler_Handle(t *testing.T) {
	validate := validator.New()
	password := &serviceMock.MockPasswordService{}
	writeModel := &modelMock.MockUserModel{}
	createHandler := handler.NewCreateHandler(validate, password, writeModel)

	ctx := context.Background()
	cmd := command.Create{
		Id: "id",
		Profile: command.CreateProfile{
			FirstName: "John",
			LastName:  "Doe",
		},
		Password: "password",
	}
	hashedPassword := domain.HashedPassword("hashedPassword")
	aUser := func(user *domain.User) bool {
		return user.GetId() == cmd.Id &&
			user.GetProfile().FirstName == cmd.Profile.FirstName &&
			user.GetProfile().LastName == cmd.Profile.LastName &&
			user.GetHashedPassword() == hashedPassword
	}

	password.On("Hash", cmd.Password).
		Return(hashedPassword, nil).
		Once()

	writeModel.On("Create", ctx, mock.MatchedBy(aUser)).
		Return(nil).
		Once()

	err := createHandler.Handle(ctx, cmd)
	assert.NoError(t, err)

	password.AssertExpectations(t)
	writeModel.AssertExpectations(t)
}
