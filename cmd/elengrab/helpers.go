package main

import (
	"fmt"
	"log"
	"os"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/pkg/nfile"
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

func ensureAssets(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot access directory %s: %w", path, err)
		}
	} else {
		return nfile.CheckDir(path)
	}

	if !assets.Embedded {
		return fmt.Errorf("directory does not exist: %s", path)
	}

	err = os.MkdirAll(path, 0o755)
	if err != nil {
		return fmt.Errorf("cannot create directory %s: %w", path, err)
	}

	err = assets.CopyToDir(path, assets.AssetsFS)
	if err != nil {
		return err
	}
	log.Printf("Embedded assets have been copied to %s\n", path)

	return nil
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
