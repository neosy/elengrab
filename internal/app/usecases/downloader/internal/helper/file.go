package helper

import (
	"os"

	"github.com/neosy/elengrab/internal/app/utils/hash"
)

func PartialHash(filePath string) *string {
	h, err := hash.FilePartialHash(filePath)
	if err == nil && h != "" {
		return &h
	}
	return nil
}

func FileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return int64(fileInfo.Size())
}
