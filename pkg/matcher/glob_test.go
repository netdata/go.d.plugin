package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobMatch_Match(t *testing.T) {
	m := globMatcher("/a/*/d")

	cases := []struct {
		expected bool
		line     string
	}{
		{
			expected: true,
			line:     "/a/b/c/d",
		},
		{
			expected: false,
			line:     "a/b/c/d",
		},
		{
			expected: false,
			line:     "This will never fail!",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, m.MatchString(c.line))
	}
}
