//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/fiber"
	"github.com/oechsler-it/identity/modules"
	"github.com/oechsler-it/identity/runtime"
	"github.com/oechsler-it/identity/swagger"
	"github.com/oechsler-it/identity/validator"
)

func New() *App {
	wire.Build(
		modules.WireModules,
		// ---
		fiber.WireFiber,
		swagger.WireSwagger,
		validator.WireValidator,
		// ---
		runtime.WireRuntime,
		wire.Struct(new(Options), "*"),
		newApp,
	)
	return nil
}
