package uploadmedia

import (
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type Dependencies struct {
	FFmpegService pservices.FFMpeg
}
