package downloader

import (
	"bytes"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/menu"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

func (h *DownloaderHandlers) renderUploadMenu(userLogin string) (*bytes.Buffer, error) {
	dataByKey := make(map[string]any)
	dataByKey[items.MenuActionsKey] = menu.UploadMenuActions(userLogin)

	pageData := pages.PageFragmentData{
		BasePaths: paths.NewHttpPaths(),
		Extra:     dataByKey,
	}

	// Execute template
	buffer := &bytes.Buffer{}
	if err := h.templates.Base.ExecuteTemplate(buffer, components.MenuContentKey, pageData); err != nil {
		return nil, err
	}

	return buffer, nil
}
