package notsimplepattern

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	sps := New()

	assert.False(t, sps.UseCache)
	assert.NotNil(t, sps.Cache)
}

func TestCreate(t *testing.T) {
	expr := "*foobar* !foo* !*bar *"

	sps, err := Create(expr)

	require.NoError(t, err)
	assert.Len(t, sps.patterns, 4)
	assert.False(t, sps.patterns[0].exclude)
	assert.True(t, sps.patterns[1].exclude)
	assert.True(t, sps.patterns[2].exclude)
	assert.False(t, sps.patterns[3].exclude)

}

func TestPatterns_Match(t *testing.T) {
	m, err := Create("*foobar* !foo* *Bar*")

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
