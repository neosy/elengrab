package avalues

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

var (
	IndexHtmlFileName                   = "index.html"
	GrabResultItemReplacemeHtmlFileName = "grab_result_item_replaceme.html"
	GrabResultItemsHistoryHtmlFileName  = "grab_result_items_history.html"
	GrabResultNewItemFirstHtmlFileName  = "grab_result_new_item_first.html"
	GrabResultNewItemHtmlFileName       = "grab_result_new_item.html"
	GrabResultWorkingHtmlFileName       = "grab_result_working.html"
	GrabResultSuccessHtmlFileName       = "grab_result_success.html"
	GrabResultStatusHtmlFileName        = "grab_result_status.html"
	GrabResultLoadHistoryHtmlFileName   = "grab_result_load_history.html"
)

func GetGrabResultStatusHtmlFileName(status dtypes.FileStatus) string {
	var name string

	switch status {
	case dtypes.FileStatusWorking:
		name = GrabResultWorkingHtmlFileName
	case dtypes.FileStatusDone:
		name = GrabResultSuccessHtmlFileName
	default:
		name = GrabResultStatusHtmlFileName
	}

	return name
}
