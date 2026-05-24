package paths

import (
	"encoding/json"
)

var (
	IndexPageJsPaths = jsScripts{
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
	}.withPath

	AuthPageJsPaths = jsScripts{
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
	}.withPath

	WatchPageJsPaths = jsScripts{
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
	}.withPath

	IndexPageJsImportJSON = jsImportFileNames{
		"utils.js",
		"constants.js",
		"cookie.js",
		"browser.js",
		"storage-state.js",
		"tooltip.js",
		"notifications.js",
		"menu.js",
		"share.js",
		"action-buttons.js",
		"player.js",
		"index.dom.js",
		"index.sse.events.js",
		"index.view.js",
		"index.menu-configs.js",
	}.jsonForTemplate

	WatchPageJsImportJSON = jsImportFileNames{
		"watch.dom.js",
	}.jsonForTemplate
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

func (scripts jsScripts) withPath(dir string) ([]JsScript, error) {
	paths, err := fileNamesWithHash(dir, scripts.fileNames())
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

func (names jsImportFileNames) nameKeysWithPath(dir string) (jsImportMap, error) {
	return assetJsNameKeyPath(names, dir)
}

func (names jsImportFileNames) jsonForTemplate(dir string) ([]byte, error) {
	jsMap, err := names.nameKeysWithPath(dir)
	if err != nil {
		return nil, err
	}

	return jsMap.JsonForTemplate()
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
