package appenv

type AppEnv string

const (
	AppEnvLocal      AppEnv = "local"
	AppEnvDevelop    AppEnv = "develop"
	AppEnvProduction AppEnv = "production"
	AppEnvTest       AppEnv = "test"
)

// String returns the string representation of the AppEnv.
func (e AppEnv) String() string {
	return string(e)
}
