package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func New() *validator.Validate {
	return validator.New()
}

var WireValidator = wire.NewSet(
	New,
)
