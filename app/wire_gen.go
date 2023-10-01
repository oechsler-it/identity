// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/oechsler-it/identity/fiber"
	"github.com/oechsler-it/identity/gorm"
	"github.com/oechsler-it/identity/modules"
	"github.com/oechsler-it/identity/modules/permission"
	app2 "github.com/oechsler-it/identity/modules/permission/app"
	fiber4 "github.com/oechsler-it/identity/modules/permission/infra/fiber"
	hook2 "github.com/oechsler-it/identity/modules/permission/infra/hook"
	model2 "github.com/oechsler-it/identity/modules/permission/infra/model"
	"github.com/oechsler-it/identity/modules/session"
	app3 "github.com/oechsler-it/identity/modules/session/app"
	fiber2 "github.com/oechsler-it/identity/modules/session/infra/fiber"
	model3 "github.com/oechsler-it/identity/modules/session/infra/model"
	"github.com/oechsler-it/identity/modules/token"
	model4 "github.com/oechsler-it/identity/modules/token/infra/model"
	"github.com/oechsler-it/identity/modules/user"
	"github.com/oechsler-it/identity/modules/user/app"
	fiber3 "github.com/oechsler-it/identity/modules/user/infra/fiber"
	"github.com/oechsler-it/identity/modules/user/infra/hook"
	"github.com/oechsler-it/identity/modules/user/infra/model"
	"github.com/oechsler-it/identity/modules/user/infra/service"
	"github.com/oechsler-it/identity/runtime"
	"github.com/oechsler-it/identity/swagger"
	"github.com/oechsler-it/identity/validator"
)

// Injectors from wire.go:

func New() *App {
	hooks := runtime.NewHooks()
	env := runtime.NewEnv()
	logger := runtime.NewLogger(env)
	runtimeRuntime := runtime.NewRuntime(hooks, logger)
	fiberApp := fiber.NewFiber(env, logger)
	quicFiber := fiber.NewQUICFiber(fiberApp)
	options := &fiber.Options{
		Env:     env,
		Hooks:   hooks,
		Logger:  logger,
		App:     fiberApp,
		QuicApp: quicFiber,
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
	gormPermissionRepo := model2.NewGormPermissionRepo(db, logger, hooks)
	findByNameHandler := app2.NewFindByNameHandler(gormPermissionRepo)
	grantPermissionHandler := app.NewGrantPermissionHandler(findByNameHandler, validate, gormUserRepo)
	createRootUser := &hook.CreateRootUser{
		Hooks:                hooks,
		Logger:               logger,
		Env:                  env,
		Repo:                 gormUserRepo,
		VerifyNoUserExists:   verifyNoUserExistsHandler,
		Create:               createHandler,
		FindPermissionByName: findByNameHandler,
		Grant:                grantPermissionHandler,
	}
	gormSessionRepo := model3.NewGormSessionRepo(db, logger, hooks)
	renewHandler := app3.NewRenewHandler(validate, gormSessionRepo)
	renewMiddleware := &fiber2.RenewMiddleware{
		Logger: logger,
		Env:    env,
		Renew:  renewHandler,
	}
	verifyActiveHandler := app3.NewVerifyActiveHandler(gormSessionRepo)
	protectMiddleware := &fiber2.ProtectMiddleware{
		VerifyActive: verifyActiveHandler,
	}
	findByIdHandler := app3.NewFindByIdHandler(gormSessionRepo)
	findByIdentifierHandler := app.NewFindByIdentifierHandler(gormUserRepo)
	userMiddleware := &fiber3.UserMiddleware{
		FindSessionById: findByIdHandler,
		FindById:        findByIdentifierHandler,
	}
	verifyHasPermissionHandler := app.NewVerifyHasPermissionHandler(gormUserRepo)
	permissionMiddleware := &fiber3.PermissionMiddleware{
		VerifyHasPermission: verifyHasPermissionHandler,
	}
	createUserHandler := &fiber3.CreateUserHandler{
		App:                  fiberApp,
		Logger:               logger,
		Validate:             validate,
		RenewMiddleware:      renewMiddleware,
		ProtectMiddleware:    protectMiddleware,
		UserMiddleware:       userMiddleware,
		PermissionMiddleware: permissionMiddleware,
		Repo:                 gormUserRepo,
		Create:               createHandler,
	}
	deleteHandler := app.NewDeleteHandler(gormUserRepo)
	deleteMeHandler := &fiber3.DeleteMeHandler{
		App:               fiberApp,
		Logger:            logger,
		RenewMiddleware:   renewMiddleware,
		ProtectMiddleware: protectMiddleware,
		UserMiddleware:    userMiddleware,
		Delete:            deleteHandler,
	}
	deleteUserHandler := &fiber3.DeleteUserHandler{
		App:                  fiberApp,
		Logger:               logger,
		RenewMiddleware:      renewMiddleware,
		ProtectMiddleware:    protectMiddleware,
		UserMiddleware:       userMiddleware,
		PermissionMiddleware: permissionMiddleware,
		Delete:               deleteHandler,
	}
	meHandler := &fiber3.MeHandler{
		App:               fiberApp,
		Logger:            logger,
		RenewMiddleware:   renewMiddleware,
		ProtectMiddleware: protectMiddleware,
		UserMiddleware:    userMiddleware,
	}
	userByIdHandler := &fiber3.UserByIdHandler{
		App:              fiberApp,
		Logger:           logger,
		FindByIdentifier: findByIdentifierHandler,
	}
	fiberGrantPermissionHandler := &fiber3.GrantPermissionHandler{
		App:                  fiberApp,
		Logger:               logger,
		RenewMiddleware:      renewMiddleware,
		ProtectMiddleware:    protectMiddleware,
		UserMiddleware:       userMiddleware,
		PermissionMiddleware: permissionMiddleware,
		Grant:                grantPermissionHandler,
	}
	revokePermissionHandler := app.NewRevokePermissionHandler(validate, gormUserRepo)
	fiberRevokePermissionHandler := &fiber3.RevokePermissionHandler{
		App:                  fiberApp,
		Logger:               logger,
		RenewMiddleware:      renewMiddleware,
		ProtectMiddleware:    protectMiddleware,
		UserMiddleware:       userMiddleware,
		PermissionMiddleware: permissionMiddleware,
		Revoke:               revokePermissionHandler,
	}
	hasPermissionHandler := &fiber3.HasPermissionHandler{
		App:    fiberApp,
		Logger: logger,
		Has:    verifyHasPermissionHandler,
	}
	userOptions := &user.Options{
		CreateRootUser:   createRootUser,
		CreateUser:       createUserHandler,
		DeleteMe:         deleteMeHandler,
		DeleteUser:       deleteUserHandler,
		Me:               meHandler,
		UserById:         userByIdHandler,
		GrantPermission:  fiberGrantPermissionHandler,
		RevokePermission: fiberRevokePermissionHandler,
		HasPermission:    hasPermissionHandler,
	}
	deviceIdMiddleware := &fiber2.DeviceIdMiddleware{
		App:    fiberApp,
		Logger: logger,
	}
	sessionIdMiddleware := &fiber2.SessionIdMiddleware{
		App: fiberApp,
	}
	initiateHandler := app3.NewInitiateHandler(validate, gormSessionRepo)
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
	revokeHandler := app3.NewRevokeHandler(validate, gormSessionRepo)
	logoutHandler := &fiber2.LogoutHandler{
		App:               fiberApp,
		Logger:            logger,
		ProtectMiddleware: protectMiddleware,
		FindById:          findByIdHandler,
		Revoke:            revokeHandler,
	}
	revokeSessionHandler := &fiber2.RevokeSessionHandler{
		App:               fiberApp,
		Logger:            logger,
		RenewMiddleware:   renewMiddleware,
		ProtectMiddleware: protectMiddleware,
		FindById:          findByIdHandler,
		Revoke:            revokeHandler,
	}
	findByOwnerUserIdHandler := app3.NewFindByOwnerUserIdHandler(gormSessionRepo)
	activeSessionsHandler := &fiber2.ActiveSessionsHandler{
		App:               fiberApp,
		RenewMiddleware:   renewMiddleware,
		ProtectMiddleware: protectMiddleware,
		FindById:          findByIdHandler,
		FindByOwnerUserId: findByOwnerUserIdHandler,
	}
	activeSessionHandler := &fiber2.ActiveSessionHandler{
		App:               fiberApp,
		RenewMiddleware:   renewMiddleware,
		ProtectMiddleware: protectMiddleware,
		FindById:          findByIdHandler,
	}
	sessionByIdHandler := &fiber2.SessionByIdHandler{
		App:               fiberApp,
		RenewMiddleware:   renewMiddleware,
		ProtectMiddleware: protectMiddleware,
		FindById:          findByIdHandler,
	}
	sessionOptions := &session.Options{
		DeviceIdMiddleware:    deviceIdMiddleware,
		SessionIdMiddleware:   sessionIdMiddleware,
		LoginHandler:          loginHandler,
		LogoutHandler:         logoutHandler,
		RevokeSessionHandler:  revokeSessionHandler,
		ActiveSessionsHandler: activeSessionsHandler,
		ActiveSessionHandler:  activeSessionHandler,
		SessionByIdHandler:    sessionByIdHandler,
	}
	verifyPermissionNotExistsHandler := app2.NewVerifyPermissionNotExistsHandler(gormPermissionRepo)
	appCreateHandler := app2.NewCreateHandler(validate, gormPermissionRepo)
	createBasePermissions := &hook2.CreateBasePermissions{
		Hooks:                     hooks,
		Logger:                    logger,
		Env:                       env,
		VerifyPermissionNotExists: verifyPermissionNotExistsHandler,
		Create:                    appCreateHandler,
	}
	fiberCreateHandler := &fiber4.CreateHandler{
		App:                  fiberApp,
		Logger:               logger,
		RenewMiddleware:      renewMiddleware,
		ProtectMiddleware:    protectMiddleware,
		UserMiddleware:       userMiddleware,
		PermissionMiddleware: permissionMiddleware,
		Create:               appCreateHandler,
	}
	appDeleteHandler := app2.NewDeleteHandler(gormPermissionRepo)
	fiberDeleteHandler := &fiber4.DeleteHandler{
		App:                  fiberApp,
		Logger:               logger,
		RenewMiddleware:      renewMiddleware,
		ProtectMiddleware:    protectMiddleware,
		UserMiddleware:       userMiddleware,
		PermissionMiddleware: permissionMiddleware,
		Delete:               appDeleteHandler,
	}
	findAllHandler := app2.NewFindAllHandler(gormPermissionRepo)
	permissionsHandler := &fiber4.PermissionsHandler{
		App:     fiberApp,
		FindAll: findAllHandler,
	}
	permissionOptions := &permission.Options{
		CreateBasePermissions: createBasePermissions,
		CreateHandler:         fiberCreateHandler,
		DeleteHandler:         fiberDeleteHandler,
		PermissionsHandler:    permissionsHandler,
	}
	gormTokenRepo := model4.NewGormTokenRepo(db, logger, hooks)
	tokenOptions := &token.Options{
		Repo: gormTokenRepo,
	}
	modulesOptions := &modules.Options{
		App:        fiberApp,
		User:       userOptions,
		Session:    sessionOptions,
		Permission: permissionOptions,
		Token:      tokenOptions,
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
