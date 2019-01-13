package matcher

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		valid        bool
		line         string
		expectedType Matcher
	}{
		{
			valid:        true,
			line:         "= full",
			expectedType: stringFullMatcher(""),
		},
		{
			valid:        true,
			line:         "~ ^prefix",
			expectedType: stringPrefixMatcher(""),
		},
		{
			valid:        true,
			line:         "~ suffix$",
			expectedType: stringSuffixMatcher(""),
		},
		{
			valid:        true,
			line:         "~ partial",
			expectedType: stringPartialMatcher(""),
		},
		{
			valid:        true,
			line:         "* glob",
			expectedType: globMatcher(""),
		},
		{
			valid:        true,
			line:         "~ regexp",
			expectedType: regexp.MustCompile(""),
		},
		{
			valid:        false,
			line:         "",
			expectedType: nil,
		},
		{
			valid:        false,
			line:         ":empty",
			expectedType: nil,
		},
		{
			valid:        true,
			line:         "!= string",
			expectedType: negMatcher{},
		},
		{
			valid:        true,
			line:         "!~ regexp",
			expectedType: negMatcher{},
		},
		{
			valid:        true,
			line:         "!* glob",
			expectedType: negMatcher{},
		},
	}
	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			m, err := Parse(test.line)
			assert.IsType(t, test.expectedType, m)
			if test.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
