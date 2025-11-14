package avalues

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func GrabResultStatusHtmlFileName(status dtypes.FileStatus) string {
	var name string

	switch status {
	case dtypes.FileStatusWorking:
		name = "grab_result_working.html"
	case dtypes.FileStatusDone:
		name = "grab_result_success.html"
	default:
		name = "grab_result_status.html"
	}

	return name
}
