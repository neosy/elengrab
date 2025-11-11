package file

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type File struct {
	logger *slog.Logger

	// repositories
	fileRep persistence.FileRepository
}

func NewFile(
	logger *slog.Logger,
	fileRep persistence.FileRepository,
) *File {
	return &File{
		logger:  logger,
		fileRep: fileRep,
	}
}
