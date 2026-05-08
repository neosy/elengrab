package uivalues

import (
	"os"
	"path/filepath"

	"github.com/neosy/elengrab/internal/pkg/httpx"
)

var (
	CssErrorFileName = "error.css"
)

var (
	CssIndexPaths = cssFileNames{
		"font-inter.css",
		"index-main.css",
		"variables.css",
		"theme-switcher.css",
		"menu.css",
		"menu-overrides.css",
		"grab-form.css",
		"result-rows.css",
		"player.css",
		"notifications.css",
	}.paths

	CssAuthPaths = cssFileNames{
		"font-inter.css",
		"auth-main.css",
		"variables.css",
		"theme-switcher.css",
	}.paths

	CssWatchPaths = cssFileNames{
		"font-inter.css",
		"watch-main.css",
		"variables.css",
		"theme-switcher.css",
	}.paths
)

type (
	cssFileNames []string
)

func (names cssFileNames) paths(dir string) ([]string, error) {
	paths, err := generateHashedFileNames(dir, names)
	if err != nil {
		return nil, err
	}

	for i, fineName := range paths {
		paths[i] = CssHttpPath(fineName)
	}

	return paths, nil
}

func CssFileRaw(fileName string, cssDir string) string {
	filePath := filepath.Join(cssDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	return httpx.MinifyCSS(string(data))
}
