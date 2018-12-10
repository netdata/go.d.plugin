package matcher

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// NG case
	_, err := New("")
	assert.Error(t, err)

	_, err = New("invalid=method")
	assert.Error(t, err)

	_, err = New("string=")
	assert.Error(t, err)

	m, err := New("string=expr")
	assert.NoError(t, err)
	assert.IsType(t, (*stringContains)(nil), m)

	// OK case
	m, err = New("string=^expr")
	assert.NoError(t, err)
	assert.IsType(t, (*stringPrefix)(nil), m)

	m, err = New("string=expr$")
	assert.NoError(t, err)
	assert.IsType(t, (*stringSuffix)(nil), m)

	m, err = New("regexp=[0-9]+")
	assert.NoError(t, err)
	assert.IsType(t, (*regexMatch)(nil), m)
}

func TestRegexMatch_Match(t *testing.T) {
	re, _ := regexp.Compile("[0-9]+")
	m := regexMatch{re}

	assert.True(t, m.Match("Match 2018"))
	assert.False(t, m.Match("No match"))

}

func TestStringContains_Match(t *testing.T) {
	m := stringContains{"apple"}

	assert.True(t, m.Match("Give me an apple, please"))
	assert.False(t, m.Match("Give me that round thing, please"))
}

func TestStringPrefix_Match(t *testing.T) {
	m := stringPrefix{"Prefix"}

	assert.True(t, m.Match("Prefix suffix"))
	assert.False(t, m.Match("Suffix prefix"))
}

func TestStringSuffix_Match(t *testing.T) {
	m := stringSuffix{"suffix"}

	assert.True(t, m.Match("Prefix suffix"))
	assert.False(t, m.Match("Suffix prefix"))
}
