package runtime

import "github.com/google/wire"

var WireRuntime = wire.NewSet(
	NewEnv,
	NewLogger,
	NewHooks,
	NewRuntime,
)
