package nfasthttp

import (
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

type ClientOption func(*fasthttp.Client)

// NewClient creates a new FastHTTP client with the specified options.
func NewClient(opts ...ClientOption) *fasthttp.Client {
	c := &fasthttp.Client{
		ReadBufferSize: readBufferSizeDefault,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// ClientOptionWithTimeout sets the read and write buffer sizes for the client
func ClientOptionWithTimeout(d time.Duration) ClientOption {
	return func(c *fasthttp.Client) {
		c.ReadTimeout = d
		c.WriteTimeout = d
	}
}

// ClientOptionWithreadBufferSize sets the read buffer size for the client
func ClientOptionWithreadBufferSize(bufferSize int) ClientOption {
	return func(c *fasthttp.Client) {
		c.ReadBufferSize = bufferSize
	}
}

func ClientOptionForceIPv4() ClientOption {
	return func(c *fasthttp.Client) {
		c.DialTimeout = func(addr string, timeout time.Duration) (net.Conn, error) {
			d := net.Dialer{
				Timeout:   timeout,
				DualStack: false,
			}
			return d.Dial("tcp4", addr)
		}
	}
}
