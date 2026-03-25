package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	iconfig "github.com/neosy/elengrab/internal/config"
	"github.com/neosy/elengrab/internal/pkg/nfile"
)

const versionFileName = "version"

// absPath resolves a relative path to an absolute path using current working directory.
// Exits the program if resolving fails.
func absPath(root, path string) string {
	path, err := nfile.AbsPath(root, path)
	if err != nil {
		log.Fatal(err.Error())
	}
	return path
}

// ensureAssets checks if the assets directory exists
func ensureAssets(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot access directory %s: %w", path, err)
		}
	} else {
		if !assets.Embedded {
			return nfile.CheckDir(path)
		}

		oldVersion := getVersionFromFile(path)
		if iconfig.AppVersion == oldVersion {
			return nfile.CheckDir(path)
		}

		var (
			assetDirName    = filepath.Base(path)
			parentAssetPath = filepath.Dir(path)
			oldAssetDirName = fmt.Sprintf("%s_%s", assetDirName, oldVersion)
		)
		if oldVersion == "" {
			oldAssetDirName = fmt.Sprintf("%s_%s_before", assetDirName, iconfig.AppVersion)
		}

		// Rename old assets directory
		oldAssetPath := filepath.Join(parentAssetPath, oldAssetDirName)
		err := os.Rename(path, oldAssetPath)
		if err != nil {
			return err
		}
		log.Printf("Rename assets to %s\n", oldAssetPath)
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

	// Create version file
	createVersionFile(path)

	return nil
}

// createVersionFile creates a version file.
func createVersionFile(dir string) error {
	var versionFilePath = filepath.Join(dir, versionFileName)

	versionFile, err := os.Create(versionFilePath)
	if err != nil {
		return err
	}
	versionFile.WriteString(iconfig.AppVersion)
	versionFile.Close()
	return nil
}

// getVersionFromFile reads the version from the version file.
func getVersionFromFile(dir string) string {
	versionFile, err := os.Open(filepath.Join(dir, versionFileName))
	if err != nil {
		return ""
	}
	defer versionFile.Close()

	versionByte, err := io.ReadAll(versionFile)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(versionByte))
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
