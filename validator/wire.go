package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

var WireValidator = wire.NewSet(
	validator.New,
)
