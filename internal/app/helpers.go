package app

import (
	"log"
	"os"

	nfile "github.com/neosy/elengrab/internal/pkg/file"
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

// ensureDirs verifies that all given directories exist and are directories.
// Returns true if all directories are valid, false otherwise.
func ensureDirs(dirs []string) bool {
	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(dir, 0o755)
				if err != nil {
					log.Printf("cannot create directory %s: %v", dir, err)
					return false
				}
				continue
			}
			log.Printf("cannot access directory %s: %v", dir, err)
			return false
		}
	}

	var allOk = true
	for _, dir := range dirs {
		if err := nfile.CheckDir(dir); err != nil {
			allOk = false
			log.Println(err)
		}
	}
	return allOk
}
