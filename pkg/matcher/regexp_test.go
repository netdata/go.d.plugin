package matcher

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpMatch_Match(t *testing.T) {
	m := RegexpMatch{regexp.MustCompile("[0-9]+")}

	assert.True(t, m.Match("Match 2018"))
	assert.False(t, m.Match("No Match"))
}
