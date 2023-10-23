package middleware

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/modules/middleware/infra/fiber"
)

var WireMiddleware = wire.NewSet(
	wire.Struct(new(fiber.AuthenticatedMiddleware), "*"),
	wire.Struct(new(fiber.AuthorizedMiddleware), "*"),
)
