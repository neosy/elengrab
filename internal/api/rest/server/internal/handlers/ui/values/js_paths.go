package uivalues

import (
	"encoding/json"
)

var (
	JsIndexPaths = jsScripts{
		{
			Path:  "htmx.min.js",
			Type:  "",
			Defer: false,
		},
		{
			Path:  "sse.min.js",
			Type:  "",
			Defer: false,
		},
		{
			Path:  "theme-switcher.js",
			Type:  "",
			Defer: true,
		},
		{
			Path:  "index-main.js",
			Type:  "module",
			Defer: false,
		},
	}.withPath

	JsAuthPaths = jsScripts{
		{
			Path:  "htmx.min.js",
			Type:  "",
			Defer: false,
		},
		{
			Path:  "sse.min.js",
			Type:  "",
			Defer: false,
		},
		{
			Path:  "theme-switcher.js",
			Type:  "",
			Defer: true,
		},
		{
			Path:  "auth-main.js",
			Type:  "module",
			Defer: false,
		},
	}.withPath

	JsIndexImportMapJSON = jsImportFileNames{
		"helper.js",
		"constants.js",
		"cookie.js",
		"tooltip.js",
		"menu.js",
		"menu-configs.js",
		"action-buttons.js",
		"index-row-event-handlers.js",
		"index-player.js",
		"index-dom-ids.js",
		"index-dom-elements.js",
	}.jsonForTemplate

	JsViewPaths = jsScripts{
		{
			Path:  "theme-switcher.js",
			Type:  "",
			Defer: true,
		},
		{
			Path:  "view-main.js",
			Type:  "module",
			Defer: false,
		},
	}.withPath
)

type (
	jsScript struct {
		Path  string
		Type  string
		Defer bool
	}

	jsScripts         []jsScript
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

func (scripts jsScripts) withPath(dir string) ([]jsScript, error) {
	paths, err := generateHashedFileNames(dir, scripts.fileNames())
	if err != nil {
		return nil, err
	}

	var newScripts = make([]jsScript, len(scripts))
	copy(newScripts, scripts)

	for i, name := range paths {
		newScripts[i].Path = JsHttpPath(name)
	}

	return newScripts, nil
}

func (names jsImportFileNames) mapWithPath(dir string) (jsImportMap, error) {
	paths, err := generateHashedFileNames(dir, names)
	if err != nil {
		return nil, err
	}

	var jsMap = make(jsImportMap, len(paths))
	for i, newName := range paths {
		key := JsHttpPath(names[i])
		jsMap[key] = JsHttpPath(newName)
	}

	return jsMap, nil
}

func (names jsImportFileNames) jsonForTemplate(dir string) ([]byte, error) {
	jsMap, err := names.mapWithPath(dir)
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
