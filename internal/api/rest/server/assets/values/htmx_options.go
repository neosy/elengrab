package avalues

import (
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	GrabResultItemHtmxOptionKey = "GrabResultItemHtmxOption"

	HtmxOptionNone                           = ""
	GrabResultItemHtmxOptionTriggerSwapOuter = `
		hx-get="{{.PathFileRow}}"
    	hx-trigger="every 1000ms"
    	hx-target="this"
    	hx-swap="outerHTML"
	`
)

func GrabResultStatusHtmxOption(status dtypes.FileStatus, hxGet string) string {
	var option string

	switch status {
	case dtypes.FileStatusNew, dtypes.FileStatusPending, dtypes.FileStatusWorking:
		option = GrabResultItemHtmxOptionTriggerSwapOuter
	default:
		option = HtmxOptionNone
	}

	option = strings.Replace(option, "{{.PathFileRow}}", hxGet, 1)

	return option
}
