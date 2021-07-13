// Package sanitize provides functions for sanitizing file paths.
package sanitize

import (
	"regexp"
	"strings"

	"github.com/alexsergivan/transliterator"
)

var (
	tl = transliterator.NewTransliterator(nil)
)

// Transliterate performs transliteration of s, returning a value suitable for use in a file or URL.
func Transliterate(s string) string {
	return tl.Transliterate(s, "")
}

// Remove all other unrecognised characters apart from
var (
	illegalName        = regexp.MustCompile(`[^[:alnum:]-._]`)
	baseNameSeparators = regexp.MustCompile(`[/]`)
)

// BaseName makes string s safe to use as a filename, producing a sanitized basename replacing . or / with -.
// NOTE: The return value may be a zero-length string.
func BaseName(s string) string {
	baseName := baseNameSeparators.ReplaceAllString(s, "-")
	return sanitize(baseName, illegalName)
}

var (
	separators = regexp.MustCompile(`[ &=+:|]`)
	dashes     = regexp.MustCompile(`[\-]+`)
)

// sanitize replaces separators with - and removes characters listed in the regexp provided from string.
// Accents, spaces, and all characters not in A-Za-z0-9 are replaced.
func sanitize(s string, r *regexp.Regexp) string {
	// Remove leading and trailing whitespace
	s = strings.Trim(s, " ")

	// Transliterate all accent characters to ASCII
	s = Transliterate(s)

	// Replace various invalid characters with a "-"
	s = separators.ReplaceAllString(s, "-")

	// Remove all other characters specified by r
	s = r.ReplaceAllString(s, "")

	// Replace multiple dashes caused by replacements above to a single dash
	s = dashes.ReplaceAllString(s, "-")

	return s
}
