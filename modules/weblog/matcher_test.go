package weblog

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMatcher(t *testing.T) {
	// NG case
	_, err := newMatcher("")
	assert.Error(t, err)

	_, err = newMatcher("invalid=method")
	assert.Error(t, err)

	_, err = newMatcher("string=")
	assert.Error(t, err)

	m, err := newMatcher("string=expr")
	assert.NoError(t, err)
	assert.IsType(t, (*stringContains)(nil), m)

	// OK case
	m, err = newMatcher("string=^expr")
	assert.NoError(t, err)
	assert.IsType(t, (*stringPrefix)(nil), m)

	m, err = newMatcher("string=expr$")
	assert.NoError(t, err)
	assert.IsType(t, (*stringSuffix)(nil), m)

	m, err = newMatcher("regexp=[0-9]+")
	assert.NoError(t, err)
	assert.IsType(t, (*regexMatch)(nil), m)
}

func TestRegexMatch_Match(t *testing.T) {
	re, _ := regexp.Compile("[0-9]+")
	m := regexMatch{re}

	assert.True(t, m.match("match 2018"))
	assert.False(t, m.match("No match"))

}

func TestStringContains_Match(t *testing.T) {
	m := stringContains{"apple"}

	assert.True(t, m.match("Give me an apple, please"))
	assert.False(t, m.match("Give me that round thing, please"))
}

func TestStringPrefix_Match(t *testing.T) {
	m := stringPrefix{"Prefix"}

	assert.True(t, m.match("Prefix suffix"))
	assert.False(t, m.match("Suffix prefix"))
}

func TestStringSuffix_Match(t *testing.T) {
	m := stringSuffix{"suffix"}

	assert.True(t, m.match("Prefix suffix"))
	assert.False(t, m.match("Suffix prefix"))
}
