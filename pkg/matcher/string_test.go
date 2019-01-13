package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringFull_Match(t *testing.T) {
	m := stringFullMatcher("Does Coca-Cola contain cocaine?")

	cases := []struct {
		expected bool
		line     string
	}{
		{
			expected: true,
			line:     "Does Coca-Cola contain cocaine?",
		},
		{
			expected: false,
			line:     "Water contains hydrogen and oxygen.",
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

func TestStringPartial_Match(t *testing.T) {
	m := stringPartialMatcher("contain")

	cases := []struct {
		expected bool
		line     string
	}{
		{
			expected: true,
			line:     "Does Coca-Cola contain cocaine?",
		},
		{
			expected: true,
			line:     "Water contains hydrogen and oxygen.",
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

func TestStringSuffix_Match(t *testing.T) {
	m := stringSuffixMatcher("mistakes.")

	cases := []struct {
		expected bool
		line     string
	}{
		{
			expected: true,
			line:     "Your paper contains too many mistakes.",
		},
		{
			expected: true,
			line:     "This sentence contains several mistakes.",
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

func TestStringPrefix_Match(t *testing.T) {
	m := stringPrefixMatcher("That book")

	cases := []struct {
		expected bool
		line     string
	}{
		{
			expected: true,
			line:     "That book contains many pictures.",
		},
		{
			expected: true,
			line:     "That book contains useful ideas.",
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
