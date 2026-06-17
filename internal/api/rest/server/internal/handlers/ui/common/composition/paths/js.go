package paths

import (
	"encoding/json"
	"errors"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
)

var (
	authPageJsPaths = jsScripts{
		{
			Path:   "htmx.min.js",
			Type:   "",
			Defer:  false,
			Legacy: LegacyNo,
		},
		{
			Path:   "htmx-1.9.12.min.js",
			Type:   "",
			Defer:  false,
			Legacy: LegacyYes,
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
			Path:   "htmx.min.js",
			Type:   "",
			Defer:  false,
			Legacy: LegacyNo,
		},
		{
			Path:   "htmx-1.9.12.min.js",
			Type:   "",
			Defer:  false,
			Legacy: LegacyYes,
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
			Path:   "htmx.min.js",
			Type:   "",
			Defer:  false,
			Legacy: LegacyNo,
		},
		{
			Path:   "htmx-1.9.12.min.js",
			Type:   "",
			Defer:  false,
			Legacy: LegacyYes,
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

	editMediaPageJsPaths = jsScripts{
		{
			Path:  "theme-switcher.js",
			Type:  "",
			Defer: true,
		},
		{
			Path:  "edit-media.page.js",
			Type:  "module",
			Defer: false,
		},
	}
)

const (
	LegacyUnknown LegacyState = iota
	LegacyYes
	LegacyNo
)

type (
	LegacyState uint8

	JsScript struct {
		Path   string
		Type   string
		Defer  bool
		Legacy LegacyState
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

func (scripts jsScripts) newLoader(assets *assets.Assets) func(legacy bool) ([]JsScript, error) {
	return func(legacy bool) ([]JsScript, error) {
		names, err := loadFileNamesWithHash(scripts.fileNames(), assetJsDir, assets)
		if err != nil {
			return nil, err
		}

		if len(names) != len(scripts) {
			return nil, errors.New("names and scripts length mismatch")
		}

		newScripts := make([]JsScript, 0, len(scripts))
		for i, script := range scripts {
			if legacy && script.Legacy == LegacyNo {
				continue
			}
			if !legacy && script.Legacy == LegacyYes {
				continue
			}

			script.Path = JsPath(names[i])
			newScripts = append(newScripts, script)
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
