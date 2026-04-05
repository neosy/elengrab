package uivalues

import iconfig "github.com/neosy/elengrab/internal/config"

const (
	pageTitle             = "Elengrab — Fast Media Downloader | YouTube, VK Video and more"
	authRegisterPageTitle = "Registration - Elengrab"
	authLoginPageTitle    = "Login - Elengrab"
	pageDescription       = "Download video and audio from YouTube, Facebook, Instagram, VK, TikTok, and more with Elengrab."
	header                = "Elengrab"
	inputGrabPlaceholder  = "Enter video or audio URL"
)

type baseValues struct {
	Title       string
	Description string

	Header     string
	AppVersion string

	ShowHeader bool
	ShowFooter bool

	AuthRegisterTitle string
	AuthLoginTitle    string
}

var (
	BaseValues = baseValues{
		Title:       pageTitle,
		Description: pageDescription,

		Header:     header,
		AppVersion: iconfig.AppVersion,

		ShowHeader: true,
		ShowFooter: true,

		AuthRegisterTitle: authRegisterPageTitle,
		AuthLoginTitle:    authLoginPageTitle,
	}

	FormGrabValues = map[string]any{
		"InputGrabPlaceholder": inputGrabPlaceholder,
	}
)

func (v baseValues) Copy() baseValues {
	return v
}

func (v baseValues) ToMap() map[string]any {
	return StructToMap(v)
}
