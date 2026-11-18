package contenttype

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type contentTypeOption struct {
	contentType     dtypes.QueryMediaContentType
	hasCreateAccess bool
}
