package broadcaster

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type clientKey struct {
	connectionID uuid.UUID
	eventKey     dtypes.EventKey
}
