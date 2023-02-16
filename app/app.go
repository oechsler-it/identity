package app

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/oechsler-it/identity/fiber"
	"github.com/oechsler-it/identity/runtime"
	"github.com/sirupsen/logrus"
)

type AppOptions struct {
	Runtime *runtime.Runtime
	Logger  *logrus.Logger
	Fiber   *fiber.FiberOptions
}

type App struct {
	AppOptions
}

func newApp(opts *AppOptions) *App {
	fiber.UseFiber(opts.Fiber)

	return &App{
		AppOptions: *opts,
	}
}

func (a *App) Run() {
	a.Runtime.Run()
}
