package assetx

import (
	"path/filepath"
	"testing"
)

func TestRewriteImports(t *testing.T) {
	fileNameWithHash := func(name string) string {
		fileName := filepath.Base(name)
		switch fileName {
		case "browser.js":
			return "browser.3994ab32.js"
		case "utils.js":
			return "utils.a10ad2c6.js"
		default:
			return name
		}
	}

	input := []byte(`
// comment before
import * as browser from "./browser.js"
/* comment 1
   comment 2
*/
/* comment */
import "./utils.js"

// comment after imports

console.log("stop here")
import "./should-not-change.js"
`)

	expected := `
// comment before
import * as browser from "./browser.3994ab32.js"
/* comment 1
   comment 2
*/
/* comment */
import "./utils.a10ad2c6.js"

// comment after imports

console.log("stop here")
import "./should-not-change.js"
`

	out := RewriteJSImports(input, fileNameWithHash)

	if string(out) != expected {
		t.Fatalf("rewriteImports failed\n--- GOT ---\n%s\n--- WANT ---\n%s\n", string(out), expected)
	}
}
