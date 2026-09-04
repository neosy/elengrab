package assets

import "html/template"

type Templates struct {
	Base  *template.Template
	Pages map[string]*template.Template
}
