package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/mahmoud24598salah/MSM_Bank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	 if currency , ok := fl.Field().Interface().(string);ok {
		return util.IsSupporredCurrency(currency)
	 }
	 return false
}