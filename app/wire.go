//go:build wireinject
// +build wireinject

//go:generate wire

package app

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/fiber"
	"github.com/oechsler-it/identity/runtime"
)

func New() *App {
	wire.Build(
		fiber.WireFiber,
		// ---
		runtime.WireRuntime,
		wire.Struct(new(AppOptions), "*"),
		newApp,
	)
	return nil
}
