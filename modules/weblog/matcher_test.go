package weblog

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newMatcher(t *testing.T) {
	// NG case
	_, err := newMatcher("")
	assert.Error(t, err)

	_, err = newMatcher("invalid=method")
	assert.Error(t, err)

	_, err = newMatcher("string=")
	assert.Error(t, err)

	m, err := newMatcher("string=expr")
	assert.NoError(t, err)
	assert.IsType(t, (*stringContainsMatcher)(nil), m)

	// OK case
	m, err = newMatcher("string=^expr")
	assert.NoError(t, err)
	assert.IsType(t, (*stringPrefixMatcher)(nil), m)

	m, err = newMatcher("string=expr$")
	assert.NoError(t, err)
	assert.IsType(t, (*stringSuffixMatcher)(nil), m)

	m, err = newMatcher("regexp=[0-9]+")
	assert.NoError(t, err)
	assert.IsType(t, (*regexpMatcher)(nil), m)
}

func Test_regexMatch_match(t *testing.T) {
	re, _ := regexp.Compile("[0-9]+")
	m := regexpMatcher{re}

	assert.True(t, m.match("MatchString 2018"))
	assert.False(t, m.match("No MatchString"))

}

func Test_stringContains_match(t *testing.T) {
	m := stringContainsMatcher{"apple"}

	assert.True(t, m.match("Give me an apple, please"))
	assert.False(t, m.match("Give me that round thing, please"))
}

func Test_stringPrefix_match(t *testing.T) {
	m := stringPrefixMatcher{"Prefix"}

	assert.True(t, m.match("Prefix suffix"))
	assert.False(t, m.match("Suffix prefix"))
}

func Test_stringSuffix_match(t *testing.T) {
	m := stringSuffixMatcher{"suffix"}

	assert.True(t, m.match("Prefix suffix"))
	assert.False(t, m.match("Suffix prefix"))
}
