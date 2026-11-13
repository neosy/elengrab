package iconfig

// AdminServerConfig contains the configuration for the admin server.
type AdminServerConfig struct {
	// Enables the admin server.
	Enable bool `env:"ENABLE" envDefault:"false"`

	// Address on which the admin server listens.
	Address string `env:"ADDRESS" envDefault:"127.0.0.1"`

	// Port on which the admin server listens.
	Port string `env:"PORT" envDefault:"6060"`

	// Debug server configuration.
	DebugConfig AdminServerDebugConfig `envPrefix:"DEBUG_"`
}

// AdminServerDebugConfig contains configuration for admin server endpoints.
type AdminServerDebugConfig struct {
	// Enables the pprof debugging endpoints.
	EnablePprof bool `env:"PPROF" envDefault:"true"`

	// Enables the metrics endpoint.
	EnableMetrics bool `env:"METRICS" envDefault:"true"`

	// Enables the health check endpoint.
	EnableHealth bool `env:"HEALTH" envDefault:"true"`
}
