package assetx

import (
	"bytes"
	"path/filepath"
)

// FileNameHasher transforms a JavaScript module file name into its
// hashed version used for cache busting (e.g. "browser.js" -> "browser.abc123.js").
type FileNameHasher func(name string) string

// RewriteJSImports rewrites static ES module import statements in a JavaScript source file,
// substituting file names with hashed versions using the provided FileNameHasher.
// Processing stops once non-import content is encountered.
func RewriteJSImports(data []byte, hashedFileName FileNameHasher) []byte {
	out := make([]byte, 0, len(data)+64)

	pos := 0
	dataLen := len(data)
	isOpenComment := false

	for pos < dataLen {
		// line start
		lineStart := pos

		// skip spaces
		for pos < dataLen && (data[pos] == ' ' || data[pos] == '\t') {
			pos++
		}

		// empty line
		if pos < dataLen && (data[pos] == '\n' || data[pos] == '\r') {
			out = append(out, data[lineStart:pos]...)
			out = append(out, data[pos])
			pos++
			continue
		}

		// find end of line
		for pos < dataLen && data[pos] != '\n' {
			pos++
		}
		line := data[lineStart:pos]

		trimmed := bytes.TrimSpace(line)

		// empty line
		if len(trimmed) == 0 {
			out = append(out, line...)
			if pos < dataLen {
				out = append(out, data[pos])
				pos++
			}
			continue
		}

		// comments
		var isComment = false
		if bytes.HasPrefix(trimmed, []byte("//")) {
			isComment = true
		}
		if bytes.HasPrefix(trimmed, []byte("/*")) {
			isComment = true
			isOpenComment = true
		}
		if isOpenComment {
			isComment = true
		}
		if bytes.HasSuffix(trimmed, []byte("*/")) {
			isComment = true
			isOpenComment = false
		}
		if isComment {
			out = append(out, line...)
			if pos < dataLen {
				out = append(out, data[pos])
				pos++
			}
			continue
		}

		// check import (very cheap, without parsing)
		if !isJSImportLine(trimmed) {
			// stop condition: do not touch the file anymore
			out = append(out, data[lineStart:]...)
			return out
		}

		// rewrite import line
		rewritten := rewriteJSImportLine(line, hashedFileName)
		out = append(out, rewritten...)

		if pos < dataLen {
			out = append(out, data[pos])
			pos++
		}
	}

	return out
}

// rewriteJSImportLine rewrites a single ES module import line by replacing
// JavaScript file names with their hashed equivalents using the provided FileNameHasher.
//
// It assumes the input line is a valid or partially valid import statement
// and does not perform full JavaScript parsing.
func rewriteJSImportLine(line []byte, hashedFileName FileNameHasher) []byte {
	out := make([]byte, 0, len(line)+16)

	pos := 0
	lineLen := len(line)

	for pos < lineLen {
		// looking for a quotation mark
		if line[pos] != '\'' && line[pos] != '"' {
			out = append(out, line[pos])
			pos++
			continue
		}

		quote := line[pos]
		out = append(out, line[pos])
		pos++

		start := pos

		// reading the path
		for pos < lineLen && line[pos] != quote {
			pos++
		}

		filePath := string(line[start:pos])

		// rewrite
		if len(filePath) > 3 && filepath.Ext(filePath) == ".js" {
			dir := filepath.Dir(filePath)
			fileNameWithHash := hashedFileName(filepath.Base(filePath))
			filePath = dir + "/" + fileNameWithHash
		}

		out = append(out, filePath...)

		if pos < lineLen {
			out = append(out, line[pos])
			pos++
		}
	}

	return out
}

func isJSImportLine(b []byte) bool {
	if len(b) < 6 {
		return false
	}
	return b[0] == 'i' &&
		b[1] == 'm' &&
		b[2] == 'p' &&
		b[3] == 'o' &&
		b[4] == 'r' &&
		b[5] == 't'
}
