package httptemplates

import (
	"fmt"
	"html/template"
	"path/filepath"
)

func LoadTemplates(path string) (*template.Template, error) {
	var (
		err  error
		tmpl = template.New("base")
	)

	tmpl, err = tmpl.ParseGlob(filepath.Join(path, "templates", "layouts", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %v", err)
	}

	tmpl, err = tmpl.ParseGlob(filepath.Join(path, "templates", "components", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %v", err)
	}

	return tmpl, nil
}
