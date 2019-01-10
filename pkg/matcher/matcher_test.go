package matcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		valid        bool
		line         string
		expectedType Matcher
	}{
		{
			valid:        true,
			line:         "=:^full$",
			expectedType: &StringFull{},
		},
		{
			valid:        true,
			line:         "=:^prefix",
			expectedType: &StringPrefix{},
		},
		{
			valid:        true,
			line:         "=:suffix$",
			expectedType: &StringSuffix{},
		},
		{
			valid:        true,
			line:         "=:partial",
			expectedType: &StringPartial{},
		},
		{
			valid:        true,
			line:         "*:glob",
			expectedType: &GlobMatch{},
		},
		{
			valid:        true,
			line:         "~:regexp",
			expectedType: &RegExpMatch{},
		},
		{
			valid:        false,
			line:         "no method",
			expectedType: nil,
		},
		{
			valid:        false,
			line:         ":empty",
			expectedType: nil,
		},
		{
			valid:        true,
			line:         "!~:regexp",
			expectedType: &NegMatcher{},
		},
		{
			valid:        true,
			line:         "!*:glob",
			expectedType: &NegMatcher{},
		},
	}

	for _, c := range cases {
		m, err := Parse(c.line)
		assert.IsType(t, c.expectedType, m)
		if c.valid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
