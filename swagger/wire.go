package swagger

import "github.com/google/wire"

var WireSwagger = wire.NewSet(
	wire.Struct(new(Options), "*"),
)
