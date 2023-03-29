// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/oechsler-it/identity/fiber"
	"github.com/oechsler-it/identity/gorm"
	"github.com/oechsler-it/identity/modules"
	"github.com/oechsler-it/identity/modules/session"
	app2 "github.com/oechsler-it/identity/modules/session/app"
	fiber2 "github.com/oechsler-it/identity/modules/session/infra/fiber"
	model2 "github.com/oechsler-it/identity/modules/session/infra/model"
	"github.com/oechsler-it/identity/modules/user"
	"github.com/oechsler-it/identity/modules/user/app"
	"github.com/oechsler-it/identity/modules/user/infra/hook"
	"github.com/oechsler-it/identity/modules/user/infra/model"
	"github.com/oechsler-it/identity/modules/user/infra/service"
	"github.com/oechsler-it/identity/runtime"
	"github.com/oechsler-it/identity/swagger"
)

// Injectors from wire.go:

func New() *App {
	hooks := runtime.NewHooks()
	env := runtime.NewEnv()
	logger := runtime.NewLogger(env)
	runtimeRuntime := runtime.NewRuntime(hooks, logger)
	fiberApp := fiber.NewFiber()
	options := &fiber.Options{
		Hooks:  hooks,
		Logger: logger,
		App:    fiberApp,
	}
	swaggerOptions := &swagger.Options{
		App: fiberApp,
	}
	gormOptions := &gorm.Options{
		Hooks:  hooks,
		Env:    env,
		Logger: logger,
	}
	db := gorm.NewPostgres(gormOptions)
	gormUserRepo := model.NewGormUserRepo(db, logger, hooks)
	verifyNoUserExistsHandler := app.NewVerifyNoUserExistsHandler(gormUserRepo)
	validate := validator.New()
	argon2idPasswordService := service.NewArgon2idPasswordService()
	createHandler := app.NewCreateHandler(validate, argon2idPasswordService, gormUserRepo)
	createRootUser := &hook.CreateRootUser{
		Hooks:              hooks,
		Logger:             logger,
		Env:                env,
		Repo:               gormUserRepo,
		VerifyNoUserExists: verifyNoUserExistsHandler,
		Create:             createHandler,
	}
	userOptions := &user.Options{
		CreateRootUser: createRootUser,
	}
	deviceIdMiddleware := &fiber2.DeviceIdMiddleware{
		App:    fiberApp,
		Logger: logger,
	}
	sessionIdMiddleware := &fiber2.SessionIdMiddleware{
		App: fiberApp,
	}
	gormSessionRepo := model2.NewGormSessionRepo(db, logger, hooks)
	initiateHandler := app2.NewInitiateHandler(validate, gormSessionRepo)
	findByIdentifierHandler := app.NewFindByIdentifierHandler(gormUserRepo)
	verifyPasswordHandler := app.NewVerifyPasswordHandler(argon2idPasswordService, gormUserRepo)
	loginHandler := &fiber2.LoginHandler{
		App:                  fiberApp,
		Logger:               logger,
		Env:                  env,
		Model:                gormSessionRepo,
		Initiate:             initiateHandler,
		FindUserByIdentifier: findByIdentifierHandler,
		VerifyPassword:       verifyPasswordHandler,
	}
	verifyActiveHandler := app2.NewVerifyActiveHandler(gormSessionRepo)
	protectMiddleware := &fiber2.ProtectMiddleware{
		VerifyActive: verifyActiveHandler,
	}
	revokeHandler := app2.NewRevokeHandler(validate, gormSessionRepo)
	logoutHandler := &fiber2.LogoutHandler{
		App:               fiberApp,
		Logger:            logger,
		ProtectMiddleware: protectMiddleware,
		Revoke:            revokeHandler,
	}
	renewHandler := app2.NewRenewHandler(validate, gormSessionRepo)
	renewMiddleware := &fiber2.RenewMiddleware{
		Logger: logger,
		Env:    env,
		Renew:  renewHandler,
	}
	findByIdHandler := app2.NewFindByIdHandler(gormSessionRepo)
	sessionHandler := &fiber2.SessionHandler{
		App:               fiberApp,
		RenewMiddleware:   renewMiddleware,
		ProtectMiddleware: protectMiddleware,
		FindById:          findByIdHandler,
	}
	sessionOptions := &session.Options{
		DeviceIdMiddleware:  deviceIdMiddleware,
		SessionIdMiddleware: sessionIdMiddleware,
		LoginHandler:        loginHandler,
		LogoutHandler:       logoutHandler,
		SessionHandler:      sessionHandler,
	}
	modulesOptions := &modules.Options{
		App:     fiberApp,
		User:    userOptions,
		Session: sessionOptions,
	}
	appOptions := &Options{
		Runtime: runtimeRuntime,
		Logger:  logger,
		Fiber:   options,
		Swagger: swaggerOptions,
		Modules: modulesOptions,
	}
	appApp := newApp(appOptions)
	return appApp
}
