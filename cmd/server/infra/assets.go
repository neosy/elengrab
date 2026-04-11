package infra

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	iconfig "github.com/neosy/elengrab/internal/config"
	nfile "github.com/neosy/elengrab/internal/pkg/file"
)

const versionFileName = "version"

func SetupAssets(cfg *iconfig.Config) error {
	assetsPath := absPath(cfg.Elengrab.RootDir, cfg.Elengrab.AssetsDir)

	if err := ensureAssets(assetsPath); err != nil {
		return fmt.Errorf("failed to ensure assets directory: %w", err)
	}

	return nil
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
