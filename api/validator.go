package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/valrichter/Ualapp/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is valid
		return util.IsValidCurrency(currency)
	}
	return false
} 