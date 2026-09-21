package mappers

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	esearchindex "github.com/neosy/elengrab/internal/repository/sqlite/search_index/entity"
)

func (m *Mappers) MapMediaSourceIndexDomainToEntity(index *ddownload.MediaSourceIndex) (*esearchindex.MediaSourceIndex, error) {
	var descriptionLower string
	if index.Description != nil {
		descriptionLower = strings.ToLower(*index.Description)
	}

	var channelID *uuid.UUID
	if index.ChannelID != uuid.Nil {
		channelID = &index.ChannelID
	}

	return &esearchindex.MediaSourceIndex{
		DownloadID:       index.DownloadID,
		UserID:           index.UserID,
		Title:            index.Title,
		TitleLower:       strings.ToLower(index.Title),
		Description:      index.Description,
		DescriptionLower: descriptionLower,
		ChannelID:        channelID,
		Visibility:       index.Visibility.String(),
		Views:            int(index.Views),
		SourceCreatedAt:  index.SourceCreatedAt,
	}, nil
}

func (m *Mappers) MapSourceIndexEntityToDomain(index *esearchindex.MediaSourceIndex) (*ddownload.MediaSourceIndex, error) {
	visibility, err := dtypes.ParseMediaVisibility(index.Visibility)
	if err != nil {
		return nil, err
	}

	var channelID uuid.UUID
	if index.ChannelID != nil {
		channelID = *index.ChannelID
	}

	return &ddownload.MediaSourceIndex{
		DownloadID:      index.DownloadID,
		UserID:          index.UserID,
		Title:           index.Title,
		Description:     index.Description,
		ChannelID:       channelID,
		Visibility:      visibility,
		Views:           uint32(index.Views),
		SourceCreatedAt: index.SourceCreatedAt,
		DeletedAt:       index.DeletedAt,
	}, nil
}

func (m *Mappers) MapRowsToSourceIndexes(rows *sql.Rows) ([]*ddownload.MediaSourceIndex, error) {
	var (
		eIndex  esearchindex.MediaSourceIndex
		indexes []*ddownload.MediaSourceIndex
	)

	for rows.Next() {
		err := rows.Scan(eIndex.FieldPointers()...)
		if err != nil {
			return nil, err
		}

		index, err := m.MapSourceIndexEntityToDomain(&eIndex)
		if err != nil {
			return nil, err
		}

		indexes = append(indexes, index)
	}

	return indexes, nil
}
