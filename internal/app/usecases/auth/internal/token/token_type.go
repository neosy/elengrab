package authtoken

type TokenType int

const (
	CookieToken TokenType = iota
	JWTToken
	APIToken
)
