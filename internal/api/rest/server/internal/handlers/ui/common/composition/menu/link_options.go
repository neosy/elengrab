package menu

import (
	"fmt"
	"html"
	"html/template"
	"regexp"
)

type linkOptions struct {
	URLTemplate string
	Placeholder string
	NewTab      bool

	DataAttributes []dataAttribute
}

type dataAttribute struct {
	Name  string
	Value string
}

func (a dataAttribute) Attribute() template.HTMLAttr {
	attr, _ := newDataAttribute(a.Name, a.Value)
	return attr
}

func newDataAttribute(name, value string) (template.HTMLAttr, error) {
	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(name) {
		return "", fmt.Errorf("invalid data attribute name %q", name)
	}

	return template.HTMLAttr(
		`data-` + name + `="` + html.EscapeString(value) + `"`,
	), nil
}
