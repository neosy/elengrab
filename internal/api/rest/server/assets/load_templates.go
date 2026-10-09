package assets

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

var tmplPaths = [][]string{
	{"templates", "layouts", "*.html"},
	{"templates", "components", "*.html"},

	{"templates", "components", "menus", "*.html"},
	{"templates", "components", "overlays", "*.html"},

	{"templates", "components", "admin", "*.html"},

	{"templates", "components", "header-content", "*.html"},
	{"templates", "components", "header-content", "blocks", "*.html"},

	{"templates", "components", "header-actions", "*.html"},
	{"templates", "components", "header-actions", "blocks", "*.html"},

	{"templates", "components", "footer-content", "*.html"},
	{"templates", "components", "footer-content", "blocks", "*.html"},

	{"templates", "components", "media-lists", "*.html"},
	{"templates", "components", "media-lists", "rows", "*.html"},
	{"templates", "components", "media-lists", "items", "*.html"},

	{"templates", "components", "detail-pages", "watch", "*.html"},
	{"templates", "components", "detail-pages", "edit-media", "*.html"},

	{"templates", "components", "dialogs", "*.html"},
	{"templates", "components", "dialogs", "content", "*.html"},
}

func LoadTemplates(assetsPath string) (*template.Template, error) {
	var (
		err  error
		tmpl = template.New("base")
	)

	tmpl = tmpl.Funcs(
		template.FuncMap{
			"lower": strings.ToLower,
		},
	)

	for _, p := range tmplPaths {
		tmpl, err = tmpl.ParseGlob(filepath.Join(append([]string{assetsPath}, p...)...))
		if err != nil {
			return nil, fmt.Errorf("error parsing templates: %v", err)
		}
	}

	return tmpl, nil
}
