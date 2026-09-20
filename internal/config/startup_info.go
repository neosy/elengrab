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

// newRuntimeInfo creates runtime information from the application configuration.
func newRuntimeInfo(c *Config) runtimeInfo {
	appMode := dtypes.MustParseAppMode(c.Elengrab.Mode)

	mediaVisibility, err := dtypes.ParseMediaVisibility(c.Elengrab.Behavior.NewMediaVisibility)
	if err != nil {
		mediaVisibility = dtypes.MediaVisibilityByAppMode(appMode)
	}

	info := runtimeInfo{
		appMode:                appMode,
		demoMode:               c.Elengrab.DemoMode,
		initialMediaVisibility: mediaVisibility,
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
