package authtoken

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generates a token based on the type:
// CookieToken -> random string for HttpOnly cookie
// JWTToken -> signed JWT with TTL
// APIToken -> random string for API use
func GenerateToken(tokenType TokenType, opts ...TokenOption) (string, error) {
	var options TokenOptions

	for _, o := range opts {
		o(&options)
	}

	switch tokenType {
	case CookieToken, APIToken:
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(b), nil

	case JWTToken:
		if options.TTL <= 0 {
			return "", errors.New("JWT TTL must be positive")
		}

		claims := jwt.MapClaims{
			"sub": options.UserID,
			"iat": time.Now().UTC().Unix(),
			"exp": time.Now().UTC().Add(options.TTL).Unix(),
			"typ": "access",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(options.JWTSecret)

	default:
		return "", errors.New("unknown token type")
	}
}
