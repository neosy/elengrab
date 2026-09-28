package mappers

import (
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
)

func (m *Mappers) MapSystemInfoToAboutDialogValues(info udto.SystemInfoResponse) pages.AboutDialogValues {
	return pages.AboutDialogValues{
		LogoIcon:   icons.FaviconIcon.FileRaw(),
		DonateIcon: icons.DonateIcon.FileRaw(),

		VersionLinkEnabled: !strings.HasSuffix(info.AppVersion, "b"),

		Version:       info.AppVersion,
		GoVersion:     info.GoVersion,
		YtDlpVersion:  info.YtDlpVersion,
		FFmpegVersion: info.FFmpegVersion,
		DenoVersion:   info.DenoVersion,
	}
}
