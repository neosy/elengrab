package infra

import (
	"log"

	"github.com/neosy/elengrab/internal/pkg/nfile"
)

// absPath resolves a relative path to an absolute path using current working directory.
// Exits the program if resolving fails.
func absPath(root, path string) string {
	path, err := nfile.AbsPath(root, path)
	if err != nil {
		log.Fatal(err.Error())
	}
	return path
}
