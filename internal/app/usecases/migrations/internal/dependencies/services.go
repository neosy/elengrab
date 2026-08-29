package dependencies

import pservices "github.com/neosy/elengrab/internal/ports/services"

type Services struct {
	Downloader pservices.Downloader
	FFMpeg     pservices.FFMpeg
}
