package dto

import "github.com/google/uuid"

type GetFileInfoResponse struct {
	YoutubeTitle         string
	FileId               uuid.UUID
	Name                 string
	Ext                  string
	FullName             string
	Path                 string
	SafeReadableFullName string
}
