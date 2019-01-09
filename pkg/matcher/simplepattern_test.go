package matcher

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellMatch_Match(t *testing.T) {
	m := ShellMatch{"*bar*"}

	cases := []struct {
		expected bool
		line     string
	}{
		{
			expected: true,
			line:     "Would you come into the bar and have a drink with me?",
		},
		{
			expected: true,
			line:     "The hotel has a licensed bar.",
		},
		{
			expected: false,
			line:     "This will never fail!",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, m.Match(c.line))
	}
}

func TestSimplePatterns_Match(t *testing.T) {
	m, err := CreateSimplePatterns("*foobar* !foo* *Bar*")

	require.NoError(t, err)

	cases := []struct {
		expected bool
		line     string
	}{
		{
			expected: false,
			line:     "Would you come into the bar and have a drink with me?",
		},
		{
			expected: true,
			line:     "His parents destined him for a career at the Bar.",
		},
		{
			expected: false,
			line:     "This will never fail!",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, m.Match(c.line))
	}
}

func TestSimplePatterns_Add(t *testing.T) {
	sps := &SimplePatterns{}

	cases := []struct {
		error   error
		pattern string
	}{
		{
			error:   nil,
			pattern: "Totally valid [pattern]",
		},
		{
			error:   nil,
			pattern: "*ally valid? [pattern]",
		},
		{
			error:   filepath.ErrBadPattern,
			pattern: "[]",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.error, sps.Add(c.pattern))
	}

	assert.Len(t, sps.Patterns, 2)
}

func TestCreateSimplePatterns(t *testing.T) {
	line := "*foobar* !foo* !*bar *"

	sps, err := CreateSimplePatterns(line)

	require.NoError(t, err)
	assert.Len(t, sps.Patterns, 4)
	assert.False(t, sps.Patterns[0].Negative)
	assert.True(t, sps.Patterns[1].Negative)
	assert.True(t, sps.Patterns[2].Negative)
	assert.False(t, sps.Patterns[3].Negative)
}
