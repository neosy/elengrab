package icons

import (
	"html/template"
	"time"

	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

const iconCacheTTL = 0 * time.Hour

var iconCache = memsimple.NewCacheWithDeaultCopier[string, iconEntry, *iconEntry]()

type iconEntry struct {
	raw template.HTML
}

func (icon *iconEntry) Copy() *iconEntry {
	return uptr.Copy(icon)
}
