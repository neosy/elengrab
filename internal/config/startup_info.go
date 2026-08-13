package iconfig

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// runtimeInfo contains application settings initialized at startup.
type runtimeInfo struct {
	appMode  dtypes.AppMode
	demoMode bool

	initialMediaVisibility dtypes.MediaVisibility
}

// startupInfo contains application settings initialized at startup.
var startupInfo runtimeInfo

// newStartupInfo creates runtime information from the application configuration.
func newStartupInfo(c *Config) runtimeInfo {
	appMode := dtypes.MustParseAppMode(c.Elengrab.Mode)

	info := runtimeInfo{
		appMode:                appMode,
		demoMode:               c.Elengrab.DemoMode,
		initialMediaVisibility: dtypes.MediaVisibilityByAppMode(appMode),
	}

	return info
}

// AppMode returns the application mode initialized at startup.
func AppMode() dtypes.AppMode {
	return startupInfo.appMode
}

// DemoMode reports whether demo mode is enabled.
func DemoMode() bool {
	return startupInfo.demoMode
}

// InitialMediaVisibility returns the initial media visibility for the application mode.
func InitialMediaVisibility() dtypes.MediaVisibility {
	return startupInfo.initialMediaVisibility
}
