package validators

import dtypes "github.com/neosy/elengrab/internal/domain/types"

func (v *Validators) register() {
	v.Validate.RegisterValidation("imageSource", dtypes.ValidateImageSource)
	v.Validate.RegisterValidation("visibility", dtypes.ValidateMediaVisibility)
}
