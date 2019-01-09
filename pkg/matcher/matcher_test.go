package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMatcher(t *testing.T) {
	cases := []struct {
		valid        bool
		line         string
		expectedType Matcher
	}{
		{
			valid:        true,
			line:         "string=hello",
			expectedType: &StringContains{},
		},
		{
			valid:        true,
			line:         "string=^hello",
			expectedType: &StringPrefix{},
		},
		{
			valid:        true,
			line:         "string=hello$",
			expectedType: &StringSuffix{},
		},
		{
			valid:        true,
			line:         "regexp=[0-9]+",
			expectedType: &RegexpMatch{},
		},
		{
			valid:        true,
			line:         "simplepattern=*foo !bar* *",
			expectedType: &SimplePatterns{},
		},
		{
			valid:        false,
			line:         "unknown=*foo !bar* *",
			expectedType: nil,
		},
		{
			valid:        false,
			line:         "no method",
			expectedType: nil,
		},
		{
			valid:        false,
			line:         "=empty",
			expectedType: nil,
		},
	}

	for _, c := range cases {
		m, err := CreateMatcher(c.line)
		assert.IsType(t, c.expectedType, m)
		if c.valid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}

}
