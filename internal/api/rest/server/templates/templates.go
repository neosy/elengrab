package httptemplates

import (
	"fmt"
	"html/template"
	"path/filepath"
)

func LoadTemplates(path string) (*template.Template, error) {
	tmpl, err := template.ParseGlob(filepath.Join(path, "templates", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %v", err)
	}

	return tmpl, nil
}
