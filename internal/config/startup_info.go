package iconfig

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type runtimeInfo struct {
	appMode  dtypes.AppMode
	demoMode bool

	initialMediaVisibility dtypes.MediaVisibility
}

var startupInfo runtimeInfo

func newStartupInfo(c *Config) runtimeInfo {
	appMode := dtypes.MustParseAppMode(c.Elengrab.Mode)

	info := runtimeInfo{
		appMode:                appMode,
		demoMode:               c.Elengrab.DemoMode,
		initialMediaVisibility: dtypes.MediaVisibilityByAppMode(appMode),
	}

	return info
}

func AppMode() dtypes.AppMode {
	return startupInfo.appMode
}

func DemoMode() bool {
	return startupInfo.demoMode
}

func InitialMediaVisibility() dtypes.MediaVisibility {
	return startupInfo.initialMediaVisibility
}
