package paths

import (
	"encoding/json"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
)

var (
	authPageJsPaths = jsScripts{
		{
			Path:  "htmx.min.js",
			Type:  "",
			Defer: false,
		},
		{
			Path:  "theme-switcher.js",
			Type:  "",
			Defer: true,
		},
		{
			Path:  "auth.page.js",
			Type:  "module",
			Defer: false,
		},
	}

	indexPageJsPaths = jsScripts{
		{
			Path:  "htmx.min.js",
			Type:  "",
			Defer: false,
		},
		{
			Path:  "theme-switcher.js",
			Type:  "",
			Defer: true,
		},
		{
			Path:  "index.page.js",
			Type:  "module",
			Defer: false,
		},
	}

	adminPageJsPaths = jsScripts{
		{
			Path:  "htmx.min.js",
			Type:  "",
			Defer: false,
		},
		{
			Path:  "theme-switcher.js",
			Type:  "",
			Defer: true,
		},
		{
			Path:  "admin.page.js",
			Type:  "module",
			Defer: false,
		},
	}

	watchPageJsPaths = jsScripts{
		{
			Path:  "theme-switcher.js",
			Type:  "",
			Defer: true,
		},
		{
			Path:  "watch.page.js",
			Type:  "module",
			Defer: false,
		},
	}
)

type (
	JsScript struct {
		Path  string
		Type  string
		Defer bool
	}

	jsScripts         []JsScript
	jsImportFileNames []string
	jsImportMap       map[string]string
)

func (scripts jsScripts) fileNames() []string {
	var names = make([]string, len(scripts))
	for i, s := range scripts {
		names[i] = s.Path
	}
	return names
}

func (scripts jsScripts) newLoader(assets *assets.Assets) func() ([]JsScript, error) {
	return func() ([]JsScript, error) {
		paths, err := loadFileNamesWithHash(scripts.fileNames(), assetJsDir, assets)
		if err != nil {
			return nil, err
		}

		newScripts := make([]JsScript, len(scripts))
		copy(newScripts, scripts)

		for i, name := range paths {
			newScripts[i].Path = JsPath(name)
		}

		return newScripts, nil
	}
}

func (m jsImportMap) JsonForTemplate() ([]byte, error) {
	return json.MarshalIndent(
		map[string]any{
			"imports": m,
		},
		"",
		"  ",
	)
}

func (names jsImportFileNames) jsonForTemplate(assets *assets.Assets, loader loaderAssetNameKeyPaths) ([]byte, error) {
	jsMap, err := names.nameKeysWithPath(assets, loader)
	if err != nil {
		return nil, err
	}

	return jsMap.JsonForTemplate()
}

func (names jsImportFileNames) nameKeysWithPath(assets *assets.Assets, loader loaderAssetNameKeyPaths) (jsImportMap, error) {
	return loader(names, assets)
}
