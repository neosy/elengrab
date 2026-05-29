package httptemplates

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

var tmplPaths = [][]string{
	{"templates", "layouts", "*.html"},
	{"templates", "components", "*.html"},
	{"templates", "components", "admin", "*.html"},
	{"templates", "components", "menus", "*.html"},
	{"templates", "components", "rows", "*.html"},
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
