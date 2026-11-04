package validators

import (
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func validateMaxBytes(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}

	maxBytes, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}

	return len(fl.Field().String()) <= maxBytes
}
