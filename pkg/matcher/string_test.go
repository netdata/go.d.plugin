package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringContains_Match(t *testing.T) {
	m := StringContains{"contains"}

	assert.True(t, m.Match("prefix contains suffix"))
	assert.False(t, m.Match("prefix suffix"))

}

func TestStringSuffix_Match(t *testing.T) {
	m := StringSuffix{"suffix"}

	assert.True(t, m.Match("Prefix suffix"))
	assert.False(t, m.Match("Suffix prefix"))
}

func TestStringPrefix_Match(t *testing.T) {
	m := StringPrefix{"Prefix"}

	assert.True(t, m.Match("Prefix suffix"))
	assert.False(t, m.Match("Suffix prefix"))
}
