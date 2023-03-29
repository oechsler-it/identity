package app

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/oechsler-it/identity/fiber"
	"github.com/oechsler-it/identity/modules"
	"github.com/oechsler-it/identity/runtime"
	"github.com/oechsler-it/identity/swagger"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Runtime *runtime.Runtime
	Logger  *logrus.Logger
	Fiber   *fiber.Options
	Swagger *swagger.Options
	// ---
	Modules *modules.Options
}

type App struct {
	Options
}

func newApp(opts *Options) *App {
	fiber.UseFiber(opts.Fiber)
	swagger.UseSwagger(opts.Swagger)
	// ---
	modules.UseModules(opts.Modules)

	return &App{
		Options: *opts,
	}
}

func (a *App) Run() {
	a.Runtime.Run()
}
