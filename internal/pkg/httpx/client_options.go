package httpx

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// ClientOption defines a function type that modifies an http.Client.
type ClientOption func(*http.Client)

// ClientOptionWithTimeout returns a ClientOption that sets the client's timeout.
// d: the duration before the client times out for a request.
func ClientOptionWithTimeout(d time.Duration) ClientOption {
	return func(c *http.Client) {
		c.Timeout = d
	}
}

// ClientOptionWithCookieJar returns a ClientOption that sets a custom cookie jar for the client.
// opt: options used to configure the cookie jar.
func ClientOptionWithCookieJar(opt *cookiejar.Options) ClientOption {
	jar, _ := cookiejar.New(opt)
	return func(c *http.Client) {
		c.Jar = jar
	}
}

// ClientOptionWithDefaultCookieJar returns a ClientOption that sets a default (empty) cookie jar for the client.
func ClientOptionWithDefaultCookieJar() ClientOption {
	jar, _ := cookiejar.New(nil)
	return func(c *http.Client) {
		c.Jar = jar
	}
}
