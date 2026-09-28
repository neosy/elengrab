package pages

import "html/template"

type AboutDialogValues struct {
	LogoIcon   template.HTML
	DonateIcon template.HTML

	VersionLinkEnabled bool

	Version       string
	GoVersion     string
	YtDlpVersion  string
	FFmpegVersion string
	DenoVersion   string
}
