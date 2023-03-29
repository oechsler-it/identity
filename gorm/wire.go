package gorm

import "github.com/google/wire"

var WireGorm = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewPostgres,
)
