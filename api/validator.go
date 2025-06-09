package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/matodrobec/simplebank/util"
)

var validCurrency validator.Func = func(filedLevel validator.FieldLevel) bool {
	if value, ok := filedLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(value)
	}

	return false
}
