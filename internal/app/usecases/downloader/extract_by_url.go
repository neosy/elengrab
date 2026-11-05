package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *downloader) ExtractByURL(
	ctx context.Context,
	authCtx dauth.AuthContext,
	url string,
) (*dto.ExtractByURLResponse, error) {
	return &dto.ExtractByURLResponse{
		Items: []dto.ExtractByURLItem{
			{
				URL:   "https://www.youtube.com/watch?v=GHLdg29hmpE",
				Title: "Hi-Fi - Не дано",
			},
		},
	}, nil
}
