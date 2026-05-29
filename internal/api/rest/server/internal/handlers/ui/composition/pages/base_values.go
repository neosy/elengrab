package pages

import (
	iconfig "github.com/neosy/elengrab/internal/config"
)

const (
	PageTitle             = "Elengrab — Fast Media Downloader | YouTube, VK Video and more"
	PageDescription       = "Download video and audio from YouTube, Facebook, Instagram, VK, TikTok, and more with Elengrab."
	PageAuthRegisterTitle = "Registration - Elengrab"
	PageAuthLoginTitle    = "Login - Elengrab"
	Header                = "Elengrab"

	IndexGrabFormInputPlaceholder = "Enter video or audio URL"
	IndexGrabGetButtonTitle       = "Get"
)

type baseValues struct {
	Title       string
	Description string

	MetaOgItems   []MetaOgItem
	MetaNameItems []MetaNameItem

	Header     string
	AppVersion string

	AuthRegisterTitle string
	AuthLoginTitle    string
}

var (
	baseValuesDefault = baseValues{
		Title:       PageTitle,
		Description: PageDescription,

		MetaOgItems:   []MetaOgItem{},
		MetaNameItems: []MetaNameItem{},

		Header:     Header,
		AppVersion: iconfig.AppVersion,
	}
)

func NewBaseValues() baseValues {
	return baseValuesDefault
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
