package nconfig

import (
	appenv "github.com/neosy/elengrab/pkg/nconfig/app_env"
	"github.com/neosy/elengrab/pkg/nlogger"
)

type AppConfig struct {
	// App mode: local, develop, test, production
	AppEnv  appenv.AppEnv `env:"APP_ENV" envDefault:"production"`
	Version string

	// values: json, console
	LogFormat nlogger.LogFormat `env:"LOG_FORMAT" envDefault:"console"`
	// values: debug, info, warn, error
	LogLevel nlogger.LogLevel `env:"LOG_LEVEL" envDefault:"warn"`
}
