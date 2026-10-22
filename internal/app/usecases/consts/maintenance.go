package consts

import "time"

// DefaultOrphanGracePeriod protects entries that are still being created:
// the index entry may appear slightly before the download row.
const DefaultOrphanGracePeriod = time.Hour
