package images

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Image struct {
	FileName string
	URL      string
	Width    int
	Height   int
}

const (
	Elengrab512ImageFileName     = "android-chrome-rounded-512x512.png"
	Elengrab1280ImageFileName    = "elengrab-1280x720.png"
	Elengrab1280ImageJpgFileName = "elengrab-1280x720.jpg"
)

func ThumbnailDefault() Image {
	return Image{
		FileName: "popcorn-bees_640x360.png",
		URL:      httppaths.BuildImagePath("popcorn-bees_640x360.png"),
		Width:    640,
		Height:   360,
	}
}

func ThumbnailVideoDefault() Image {
	return Image{
		FileName: "video-film_640x360.png",
		URL:      httppaths.BuildImagePath("video-film_640x360.png"),
		Width:    640,
		Height:   360,
	}
}

func ThumbnailMusicDefault() Image {
	return Image{
		FileName: "music-wave_640x360.png",
		URL:      httppaths.BuildImagePath("music-wave_640x360.png"),
		Width:    640,
		Height:   360,
	}
}

func (i Image) ImageData() *dtypes.ImageData {
	return &dtypes.ImageData{
		URL:    i.URL,
		Format: dtypes.ImageFormatFromFileName(i.FileName),
		Width:  i.Width,
		Height: i.Height,
	}
}
