package iconfig

// HTTPServerConfig contains the HTTP server configuration.
type HTTPServerConfig struct {
	// Address specifies the network address on which the HTTP server listens.
	Address string `env:"ADDRESS" envDefault:""`

	// Port specifies the port on which the HTTP server listens.
	Port string `env:"PORT" envDefault:"8080"`

	// Compress enables HTTP response compression.
	Compress bool `env:"COMPRESS" envDefault:"true"`

	// Reverse proxy authentication settings.
	ReverseProxy HTTPServerReverseProxyConfig `envPrefix:"REVERSE_PROXY_"`
}

// HTTPServerReverseProxyConfig contains the reverse proxy configuration for the HTTP server.
type HTTPServerReverseProxyConfig struct {
	// Enables authentication using the username provided in the X-Auth-Request-User header by trusted reverse proxies.
	EnableAuth bool `env:"ENABLE_AUTH" envDefault:"false"`

	// Specifies trusted reverse proxy IP addresses and CIDR ranges (e.g. "127.0.0.1/32,172.18.0.0/16").
	TrustedProxies []string `env:"TRUSTED_PROXIES" envDefault:""`
}
