package authweb

import (
	"regexp"

	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

var (
	loginRegexp    = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{4,19}$`)
	passwordRegexp = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*]{8,}$`)
	letterRegexp   = regexp.MustCompile(`[A-Za-z]`)
	digitRegexp    = regexp.MustCompile(`\d`)
	specialRegexp  = regexp.MustCompile(`[!@#$%^&*]`)
)

func isValidLogin(login string) bool {
	return loginRegexp.MatchString(login)
}

func isValidPassword(password string) bool {
	if !passwordRegexp.MatchString(password) {
		return false
	}
	if !letterRegexp.MatchString(password) {
		return false
	}
	if !digitRegexp.MatchString(password) {
		return false
	}
	if !specialRegexp.MatchString(password) {
		return false
	}
	return true
}

func checkValidLogin(login string) error {
	if !isValidLogin(login) {
		return errorx.NewHTTPMessage("login must start with a letter, 5-20 letters/digits/_", fasthttp.StatusBadRequest)
	}
	return nil
}

func checkValidPassword(login string) error {
	if !isValidPassword(login) {
		return errorx.NewHTTPMessage("password must be at least 8 chars, include a letter, a number and a special symbol (!@#$%^&*)", fasthttp.StatusBadRequest)
	}
	return nil
}
