package filex

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	NameMaxBytes     = 255 // typical FS limit
	nameMaxExtBytes  = 16  // longer "extensions" are treated as part of the name
	httpNameFallback = "file"
)

var (
	// Windows forbids: \ / : * ? " < > |
	// Unix forbids: /
	fsNameSeparatorsRe  = regexp.MustCompile(`[:|/\\]`)
	fsNameRemovableRe   = regexp.MustCompile(`["<>]`)
	fsNameReplaceableRe = regexp.MustCompile(`[?*\x00-\x1F]`)
	fsNameUnderscoresRe = regexp.MustCompile(`_+`)

	// '#' and '%' break URLs (fragment and percent-encoding).
	httpNameReplacer          = strings.NewReplacer("#", "--", "%", "_")
	httpNameMultiSpaceRe      = regexp.MustCompile(`\s+`)
	httpNameMultiUnderscoreRe = regexp.MustCompile(`_+`)
	httpNameSpaceBeforeExtRe  = regexp.MustCompile(`\s+\.`)
)

// SanitizeName replaces characters not allowed in file names on common
// filesystems (Windows and Unix) and truncates the result to 255 bytes
// (UTF-8 safe), preserving a short extension.
// The result may be empty if the input contains nothing usable.
func SanitizeName(name string) string {
	safe := strings.TrimSpace(name)

	safe = fsNameSeparatorsRe.ReplaceAllString(safe, "-")
	safe = fsNameRemovableRe.ReplaceAllString(safe, "")
	safe = fsNameReplaceableRe.ReplaceAllString(safe, "_")
	safe = fsNameUnderscoresRe.ReplaceAllString(safe, "_")

	return truncateNameKeepExt(safe, NameMaxBytes)
}

// SanitizeHTTPName is like SanitizeName but also makes the name safe for URLs
// and HTTP responses (sendfile, Content-Disposition): it replaces '#' and '%',
// normalizes whitespace and never returns an empty name ("file" is used instead).
func SanitizeHTTPName(name string) string {
	safe := SanitizeName(name)

	safe = httpNameReplacer.Replace(safe)
	safe = httpNameMultiSpaceRe.ReplaceAllString(safe, " ")
	safe = httpNameMultiUnderscoreRe.ReplaceAllString(safe, "_")
	safe = httpNameSpaceBeforeExtRe.ReplaceAllString(safe, ".")

	// '#' -> '--' makes the name longer, so truncate again.
	safe = truncateNameKeepExt(safe, NameMaxBytes)
	safe = strings.TrimSpace(safe)

	if safe == "" || safe == "." || safe == ".." {
		return httpNameFallback
	}

	return safe
}

// truncateNameKeepExt truncates a file name to at most maxBytes bytes,
// cutting the base name and keeping a short extension intact.
func truncateNameKeepExt(name string, maxBytes int) string {
	if len(name) <= maxBytes {
		return name
	}

	ext := filepath.Ext(name)
	if ext == "" || len(ext) > nameMaxExtBytes || len(ext) >= maxBytes {
		return truncateBytesUTF8(name, maxBytes)
	}

	base := truncateBytesUTF8(strings.TrimSuffix(name, ext), maxBytes-len(ext))
	return base + ext
}

// truncateBytesUTF8 truncates s to at most maxBytes bytes without splitting a rune.
func truncateBytesUTF8(s string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(s) <= maxBytes {
		return s
	}

	for maxBytes > 0 && !utf8.RuneStart(s[maxBytes]) {
		maxBytes--
	}

	return s[:maxBytes]
}
