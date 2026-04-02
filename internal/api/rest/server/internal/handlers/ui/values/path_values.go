package uivalues

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

const (
	PathFileRowKey      = "PathFileRow"
	PathFilesHistoryKey = "PathFilesHistory"
	PathStaticImgKey    = "PathStaticImg"
	PathStaticIconsKey  = "PathStaticIcons"
	PathStaticCssKey    = "PathStaticCss"
	PathStaticJsKey     = "PathStaticJs"
	PathStaticPwaKey    = "PathStaticPwa"
	PathItemsHistoryKey = "PathItemsHistory"

	PathAuthRegisterKey = "PathAuthRegister"
	PathAuthLoginKey    = "PathAuthLogin"

	PathDownloaderKey     = "PathDownloader"
	PathAccountMenuKey    = "PathAccountMenu"
	PathDownloaderGrabKey = "PathDownloaderGrab"
	PathHistorySearchKey  = "PathHistorySearch"
)

var PathValues = map[string]any{
	PathStaticImgKey:   httppaths.GroupStatic + "/img",
	PathStaticIconsKey: httppaths.GroupStatic + "/img/icons",
	PathStaticCssKey:   httppaths.GroupStatic + "/css",
	PathStaticJsKey:    httppaths.GroupStatic + "/js",
	PathStaticPwaKey:   httppaths.GroupStatic + "/pwa",

	PathAuthRegisterKey: httppaths.GroupAccount + httppaths.PathRegister,
	PathAuthLoginKey:    httppaths.GroupAccount + httppaths.PathLogin,

	PathDownloaderKey:     httppaths.GroupDownloader,
	PathAccountMenuKey:    httppaths.GroupDownloader + httppaths.PathAccountMenu,
	PathItemsHistoryKey:   httppaths.GroupDownloader + httppaths.PathHistory,
	PathDownloaderGrabKey: httppaths.GroupDownloader + httppaths.PathGrab,
	PathHistorySearchKey:  httppaths.GroupDownloader + httppaths.PathSearch,
}

var cssFileNames = []string{
	"index.css",
	"variables.css",
	"theme-switcher.css",
	"account-menu.css",
	"grab-form.css",
	"result-rows.css",
	"player.css",
}

var cssAuthFileNames = []string{
	"auth-main.css",
	"variables.css",
	"theme-switcher.css",
}

var jsImportFileNames = []string{
	"helper.js",
	"cookie.js",
	"action-button.js",
	"row-event-handlers.js",
	"player.js",
	"tooltip.js",
	"account-menu.js",
	"dom-ids.js",
	"dom-elements.js",
	"constants.js",
}

var jsScripts = []JsScript{
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
		Path:  "main.js",
		Type:  "module",
		Defer: false,
	},
}

var jsAuthScripts = []JsScript{
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
}

var (
	jsFileNames     = make([]string, len(jsScripts))
	jsAuthFileNames = make([]string, len(jsAuthScripts))
)

func init() {
	for i, s := range jsScripts {
		jsFileNames[i] = s.Path
	}

	for i, s := range jsAuthScripts {
		jsAuthFileNames[i] = s.Path
	}
}

func generateHashedFileNames(dir string, fileNames []string) ([]string, error) {
	var result = make([]string, len(fileNames))

	for i, f := range fileNames {
		filePath := filepath.Join(dir, f)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		hash := fmt.Sprintf("%x", sha256.Sum256(data))[:8] // короткий хэш
		ext := filepath.Ext(f)
		name := strings.TrimSuffix(f, ext)
		result[i] = fmt.Sprintf("%s.%s%s", name, hash, ext)
	}

	return result, nil
}

func CssPaths(dir string) ([]string, error) {
	paths, err := generateHashedFileNames(dir, cssFileNames)
	if err != nil {
		return nil, err
	}

	prefix := httppaths.GroupStatic + "/css/"

	for i, path := range paths {
		paths[i] = prefix + path
	}

	return paths, nil
}

func CssAuthPaths(dir string) ([]string, error) {
	paths, err := generateHashedFileNames(dir, cssAuthFileNames)
	if err != nil {
		return nil, err
	}

	prefix := httppaths.GroupStatic + "/css/"

	for i, path := range paths {
		paths[i] = prefix + path
	}

	return paths, nil
}

type JsScript struct {
	Path  string
	Type  string
	Defer bool
}

func JsScripts(dir string) ([]JsScript, error) {
	paths, err := generateHashedFileNames(dir, jsFileNames)
	if err != nil {
		return nil, err
	}

	prefix := httppaths.GroupStatic + "/js/"

	var scripts = make([]JsScript, len(jsScripts))
	copy(scripts, jsScripts)

	for i, path := range paths {
		scripts[i].Path = prefix + path
	}

	return scripts, nil
}

func JsAuthScripts(dir string) ([]JsScript, error) {
	paths, err := generateHashedFileNames(dir, jsAuthFileNames)
	if err != nil {
		return nil, err
	}

	prefix := httppaths.GroupStatic + "/js/"

	var scripts = make([]JsScript, len(jsAuthScripts))
	copy(scripts, jsAuthScripts)

	for i, path := range paths {
		scripts[i].Path = prefix + path
	}

	return scripts, nil
}

func JsImportMap(dir string) (map[string]string, error) {
	paths, err := generateHashedFileNames(dir, jsImportFileNames)
	if err != nil {
		return nil, err
	}

	prefix := httppaths.GroupStatic + "/js/"

	var jsMap = make(map[string]string, len(paths))
	for i, path := range paths {
		key := prefix + jsImportFileNames[i]
		jsMap[key] = prefix + path
	}

	return jsMap, nil
}
