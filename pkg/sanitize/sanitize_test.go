package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseName(t *testing.T) {
	tests := []struct {
		name string
		s    string
		exp  string
	}{
		{
			name: "no change",
			s:    "foo",
			exp:  "foo",
		},
		{
			name: "trims spaces",
			s:    " foo ",
			exp:  "foo",
		},
		{
			name: "transliterates",
			s:    "fôöd_åœ_®©éñ.jpg",
			exp:  "food_aoe_rcen.jpg",
		},
		{
			name: "replaces special",
			s:    "foo/bar|foo&bar+foo=bar:foo",
			exp:  "foo-bar-foo-bar-foo-bar-foo",
		},
		{
			name: "replaces other chars",
			s:    "foo\nbar\bfoo\tbar\afoo\r",
			exp:  "foobarfoobarfoo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BaseName(tt.s)
			assert.Equal(t, tt.exp, got)
		})
	}
}
