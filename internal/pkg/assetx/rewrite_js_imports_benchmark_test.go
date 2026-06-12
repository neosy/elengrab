package assetx

import (
	"path/filepath"
	"testing"
)

func BenchmarkRewriteImports(b *testing.B) {
	var fileNameWithHash = func(name string) string {
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

	var benchInput = []byte(`
// comment before
import * as browser from "./browser.js"
import "./utils.js"
import "./browser.js"
import "./utils.js"
import "./browser.js"
import "./utils.js"
import "./browser.js"
import "./utils.js"

// comment after imports

console.log("stop here")
import "./should-not-change.js"
`)

	b.ReportAllocs()

	for b.Loop() {
		_ = RewriteJSImports(benchInput, fileNameWithHash)
	}
}
