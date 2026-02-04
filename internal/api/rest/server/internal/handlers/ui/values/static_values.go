package uivalues

import iconfig "github.com/neosy/elengrab/infrastructure/config"

const (
	pageTitle            = "Elengrab — Fast Media Downloader | YouTube, VK Video and more"
	pageDescription      = "Download video and audio from YouTube, Facebook, Instagram, VK, TikTok, and more with Elengrab."
	header               = "Elengrab"
	inputGrabPlaceholder = "Enter video or audio URL"
)

var (
	IndexValues = map[string]any{
		"Title":       pageTitle,
		"Description": pageDescription,
		"Header":      header,
		"AppVersion":  iconfig.AppVersion,
	}

	FormGrabValues = map[string]any{
		"InputGrabPlaceholder": inputGrabPlaceholder,
	}
)
