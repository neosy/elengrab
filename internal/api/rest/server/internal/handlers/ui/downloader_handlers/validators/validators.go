package validators

import (
	"github.com/go-playground/validator/v10"
)

type Validators struct {
	Validate *validator.Validate
}

func NewValidators() *Validators {
	validators := &Validators{
		Validate: validator.New(),
	}

	validators.register()

	return validators
}
