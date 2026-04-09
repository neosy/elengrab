package httpx

import (
	"bytes"
	"log"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

// MinifyCSS minifies the given CSS string and returns the minified version.
// If an error occurs during minification, it logs the error and returns the original input.
func MinifyCSS(input string) string {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)

	var out bytes.Buffer
	err := m.Minify("text/css", &out, strings.NewReader(input))
	if err != nil {
		log.Println("minify error:", err)
		return input
	}

	return out.String()
}
