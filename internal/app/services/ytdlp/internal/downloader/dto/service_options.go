package idto

type ServiceOptions struct {
	// Number of fragments of a dash/hlsnative video that should be downloaded concurrently (default is 1)
	ConcurrentFragments uint8

	// CookiesDir path for storing cookies
	CookiesDir string
	// AllowCookies allow cookies for YouTube, Instagram, etc
	AllowCookies bool
}
