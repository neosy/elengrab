package menu

import "html/template"

type menuAction struct {
	RenderType renderType    `json:"renderType"` // "link" | "action" | "divider"
	Action     string        `json:"action"`
	Title      string        `json:"title"`
	IconSvg    template.HTML `json:"iconSvg,omitempty"`
	Link       linkOptions
	Text       string
}
