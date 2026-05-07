package fsstorage

import (
	"fmt"

	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Storages struct {
	Thumbnail pstorage.ThumbnailsStorage
	Download  pstorage.DownloadsStorage
}

func NewStorages(
	thumbnailPath string,
	downloadPath string,
	mediaDirName string,
) (*Storages, error) {
	thumbnail, err := newThumbnailsStorage(thumbnailPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize thumbnail Storage: %w", err)
	}

	download, err := newDownloadsStorage(downloadPath, mediaDirName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize download Storage: %w", err)
	}

	return &Storages{
		Thumbnail: thumbnail,
		Download:  download,
	}, nil
}
