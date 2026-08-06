package downloadpreparer

import (
	"context"
	"errors"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

type fetchMediaInfoResult struct {
	info        *idto.ExtractInfo
	mediaFormat *idto.ExtractMediaFormat

	selectedFormatIDs []string
}

func (p *DownloadPreparer) fetchMediaInfo(
	ctx context.Context,
	url string,
	formatQuery string,
	dlOptions idto.DLOptions,
) (*fetchMediaInfoResult, error) {
	// Get information about the best format from yt-dlp
	var err error
	info, err := p.getExtractInfo(
		ctx,
		url,
		formatQuery,
		idto.WithUseCookies(dlOptions.CookieFilePathIfNeeded()),
		idto.WithEnsureCache(false),
	)
	if err != nil {
		return nil, err
	}

	if len(info.Formats) == 0 {
		return nil, errors.New("not found format")
	}

	mediaFormat := info.Formats[0]

	formatIDs := append([]string{}, mediaFormat.FormatID)

	if len(info.Formats) == 2 && info.Formats[1].ACodec != "" {
		mediaFormat.ACodec = info.Formats[1].ACodec
		mediaFormat.Abr = info.Formats[1].Abr
		mediaFormat.Asr = info.Formats[1].Asr

		formatIDs = append(formatIDs, info.Formats[1].FormatID)
	}

	return &fetchMediaInfoResult{
		info:              info,
		mediaFormat:       &mediaFormat,
		selectedFormatIDs: formatIDs,
	}, nil
}
