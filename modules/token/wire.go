package token

import (
	"github.com/google/wire"
	"github.com/oechsler-it/identity/modules/token/infra/model"
)

type Options struct {
	Repo *model.GormTokenRepo
}

func UseToken(opts *Options) {
}

var WireToken = wire.NewSet(
	wire.Struct(new(Options), "*"),

	model.NewGormTokenRepo,
)
