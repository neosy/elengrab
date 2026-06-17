package pages

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

type (
	PagePaths struct {
		Css []string

		JsScripts       []paths.JsScript
		JsImportMapJSON template.HTML

		PwaManifest string
		StreamURL   string
	}
)
