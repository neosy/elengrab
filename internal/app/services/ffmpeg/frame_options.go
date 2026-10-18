package ffmpegsrv

import (
	"github.com/neosy/elengrab/internal/app/services/ffmpeg/ffframe"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func newFrameOptionsDefault() pservices.FrameOptions {
	return pservices.FrameOptions{
		Strategy: ffframe.FrameStrategyThumbnail{},
		Format:   ffframe.FrameFormatJPEG{},
	}
}
