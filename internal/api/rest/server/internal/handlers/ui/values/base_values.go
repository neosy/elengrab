package uivalues

import (
	iconfig "github.com/neosy/elengrab/internal/config"
)

const (
	PageTitle             = "Elengrab — Fast Media Downloader | YouTube, VK Video and more"
	PageDescription       = "Download video and audio from YouTube, Facebook, Instagram, VK, TikTok, and more with Elengrab."
	PageAuthRegisterTitle = "Registration - Elengrab"
	PageAuthLoginTitle    = "Login - Elengrab"
	header                = "Elengrab"
	inputGrabPlaceholder  = "Enter video or audio URL"
)

type baseValues struct {
	Title       string
	Description string

	MetaOgItems   []MetaOgItem
	MetaNameItems []MetaNameItem

	Header     string
	AppVersion string

	ShowHeader bool
	ShowFooter bool

	AuthRegisterTitle string
	AuthLoginTitle    string
}

type errorValues struct {
	Title  string
	Header string

	ErrorCode      int
	ErrorTitle     string
	ErrorText      string
	DebugErrorText any
}

var (
	defaultBaseValues = baseValues{
		Title:       PageTitle,
		Description: PageDescription,

		MetaOgItems:   make([]MetaOgItem, 0),
		MetaNameItems: make([]MetaNameItem, 0),

		Header:     header,
		AppVersion: iconfig.AppVersion,

		ShowHeader: true,
		ShowFooter: true,
	}

	defaultErrorValues = errorValues{
		Header: header,
	}

	FormGrabValues = map[string]any{
		"InputGrabPlaceholder": inputGrabPlaceholder,
	}
)

func NewBaseValues() baseValues {
	return defaultBaseValues
}

func NewErrorValues() errorValues {
	return defaultErrorValues
}

func (v baseValues) ToMap() map[string]any {
	return StructToMap(v)
}

func (v errorValues) ToMap() map[string]any {
	return StructToMap(v)
}

type (
	MetaOgItem struct {
		OgProperty string
		Content    string
	}
	MetaNameItem struct {
		Name    string
		Content string
	}
	MetaOgItems   []MetaOgItem
	MetaNameItems []MetaNameItem
)

func (m *MetaOgItems) Add(ogProperty, content string) {
	if m == nil {
		return
	}

	*m = append(*m, MetaOgItem{
		OgProperty: ogProperty,
		Content:    content,
	})
}

func (m *MetaNameItems) Add(name, content string) {
	if m == nil {
		return
	}

	*m = append(*m, MetaNameItem{
		Name:    name,
		Content: content,
	})
}
