package uivalues

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

var (
	CssIndexPaths = cssFileNames{
		"index.css",
		"variables.css",
		"theme-switcher.css",
		"menu.css",
		"grab-form.css",
		"result-rows.css",
		"player.css",
	}.paths

	CssAuthPaths = cssFileNames{
		"auth-main.css",
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
		paths[i] = httppaths.GroupCss + "/" + fineName
	}

	return paths, nil
}
