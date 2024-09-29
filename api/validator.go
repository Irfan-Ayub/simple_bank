package api

import (
	"github.com/Irfan-Ayub/simple_bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is suported
		return util.IsSupportedCurrency(currency)
	}

	return false
}
