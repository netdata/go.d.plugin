package matcher

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegExpMatch_Match(t *testing.T) {
	m := regexp.MustCompile("[0-9]+")

	cases := []struct {
		expected bool
		line     string
	}{
		{
			expected: true,
			line:     "2019",
		},
		{
			expected: true,
			line:     "It's over 9000!",
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
