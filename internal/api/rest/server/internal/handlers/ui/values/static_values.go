package uivalues

import iconfig "github.com/neosy/elengrab/infrastructure/config"

const (
	pageTitle            = "Elengrab — Fast YouTube Downloader"
	header               = "Elengrab"
	inputGrabPlaceholder = "Enter video or audio URL"
)

var (
	IndexValues = map[string]any{
		"Title":      pageTitle,
		"Header":     header,
		"AppVersion": iconfig.AppVersion,
	}

	FormGrabValues = map[string]any{
		"InputGrabPlaceholder": inputGrabPlaceholder,
	}
)
