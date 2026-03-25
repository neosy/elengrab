package authtoken

import "time"

type TokenOption func(*TokenOptions)

type TokenOptions struct {
	UserID    string
	TTL       time.Duration
	JWTSecret []byte
}
