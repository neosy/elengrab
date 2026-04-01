package dtypes

import "strings"

type Login string

func NewLogin(raw string) Login {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	return Login(normalized)
}

func (l Login) String() string {
	return string(l)
}
