package page

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
)

type Page struct {
	Name       PageName
	Title      string
	LogoSymbol string
	LogoIcon   icons.Icon

	ContentTemplateName string
}

func (p *Page) HasIcon() bool {
	return p.LogoIcon.FileName() != ""
}

func (p *Page) LogoSymbolHTML() template.HTML {
	return template.HTML(p.LogoSymbol)
}

func (p *Page) LogoIconHTML() template.HTML {
	return p.LogoIcon.FileRaw()
}

func (p *Page) LogoHTML() template.HTML {
	if p.HasIcon() {
		return p.LogoIconHTML()
	}
	return p.LogoSymbolHTML()
}
