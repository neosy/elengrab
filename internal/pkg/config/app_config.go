package nconfig

import (
	appenv "github.com/neosy/elengrab/internal/pkg/config/app_env"
	"github.com/neosy/elengrab/internal/pkg/logger"
)

type AppConfig struct {
	// AppEnv defines the application mode: local, develop, test, or production.
	AppEnv appenv.AppEnv `env:"APP_ENV" envDefault:"production"`

	// LogFormat specifies the output format for logs: "json" or "console".
	LogFormat nlogger.LogFormat `env:"LOG_FORMAT" envDefault:"console"`
	// LogLevel sets the minimum log level: "debug", "info", "warn", or "error".
	LogLevel nlogger.LogLevel `env:"LOG_LEVEL" envDefault:"warn"`
}
