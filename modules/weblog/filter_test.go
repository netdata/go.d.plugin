package weblog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawFilter_String(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		include  string
		exclude  string
	}{
		{"empty", `{"include": "", "exclude": ""}`, "", ""},
		{"include", `{"include": "a", "exclude": ""}`, "a", ""},
		{"exclude", `{"include": "", "exclude": "b"}`, "", "b"},
		{"both", `{"include": "a", "exclude": "b\"c"}`, "a", `b"c`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, rawFilter{Include: test.include, Exclude: test.exclude}.String())
		})
	}
}

func TestNewFilter(t *testing.T) {
	lines := []string{
		"",
		"abc",
		"abz",
		"bbc",
		"bbz",
	}
	tests := []struct {
		name    string
		include string
		exclude string
		answer  []bool
	}{
		{"empty", "", "", []bool{true, true, true, true, true}},
		{"include", "~ ^a", "", []bool{false, true, true, false, false}},
		{"exclude", "", "~ z$", []bool{true, true, false, true, false}},
		{"both", "~ ^a", "~ z$", []bool{false, true, false, false, false}},
	}
	for _, test := range tests {
		filter, err := NewFilter(rawFilter{
			Include: test.include,
			Exclude: test.exclude,
		})
		assert.NoError(t, err)
		assert.NotNil(t, filter)
		for i, line := range lines {
			t.Run(test.name+"_"+line, func(t *testing.T) {
				assert.Equal(t, test.answer[i], filter.MatchString(line))
			})
		}
	}
}
