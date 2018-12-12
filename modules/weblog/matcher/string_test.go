package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringContains_Match(t *testing.T) {
	m := StringContainsMatcher{"apple"}

	assert.True(t, m.Match("Give me an apple, please"))
	assert.False(t, m.Match("Give me that round thing, please"))
}

func TestStringPrefix_Match(t *testing.T) {
	m := StringPrefixMatcher{"Prefix"}

	assert.True(t, m.Match("Prefix suffix"))
	assert.False(t, m.Match("Suffix prefix"))
}

func TestStringSuffix_Match(t *testing.T) {
	m := StringSuffixMatcher{"suffix"}

	assert.True(t, m.Match("Prefix suffix"))
	assert.False(t, m.Match("Suffix prefix"))
}
