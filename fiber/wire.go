package fiber

import "github.com/google/wire"

var WireFiber = wire.NewSet(
	wire.Struct(new(FiberOptions), "*"),
	NewFiber,
)
