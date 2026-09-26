package pagesdata

import (
	"strconv"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/images"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

func BuildMetaOgItems(baseURL string) pages.MetaOgItems {
	imageData := &dtypes.ImageData{
		URL:    baseURL + paths.ImagePath(images.Elengrab1280ImageJpgFileName),
		Format: dtypes.ImageFormatJPEG,
		Width:  1280,
		Height: 720,
	}

	metaOgItems := make(pages.MetaOgItems, 0, 15)
	metaOgItems.Add("site_name", iconfig.AppName)
	metaOgItems.Add("type", "website")
	metaOgItems.Add("title", pages.PageTitle)
	metaOgItems.Add("description", pages.PageDescription)
	metaOgItems.Add("url", baseURL)
	metaOgItems.Add("image", imageData.URL)
	metaOgItems.Add("image:secure_url", imageData.URL)
	metaOgItems.Add("image:type", httpx.ContentTypeByExt(imageData.Format.String()))
	metaOgItems.Add("image:width", strconv.Itoa(imageData.Width))
	metaOgItems.Add("image:height", strconv.Itoa(imageData.Height))
	metaOgItems.Add("image:alt", "Elengrab logo")

	return metaOgItems
}
