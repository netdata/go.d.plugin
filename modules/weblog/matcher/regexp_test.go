package matcher

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexMatch_Match(t *testing.T) {
	re, _ := regexp.Compile("[0-9]+")
	m := RegexpMatcher{re}

	assert.True(t, m.Match("Match 2018"))
	assert.False(t, m.Match("No Match"))

}
