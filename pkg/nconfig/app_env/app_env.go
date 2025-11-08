package appenv

type AppEnv string

const (
	AppEnvLocal      AppEnv = "local"
	AppEnvDevelop    AppEnv = "develop"
	AppEnvProduction AppEnv = "production"
	AppEnvTest       AppEnv = "test"
)
