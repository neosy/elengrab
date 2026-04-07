package uivalues

var (
	CssIndexPaths = cssFileNames{
		"index-main.css",
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

	CssViewPaths = cssFileNames{
		"view-main.css",
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
