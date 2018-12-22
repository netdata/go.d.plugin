package csvparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_ParseString(t *testing.T) {
	tests := []struct {
		name   string
		expect []string
		err    bool
	}{
		{"", nil, false},
		{",", []string{"", ""}, false},
		{"a,b,", []string{"a", "b", ""}, false},
		{"a,b\nc,d", []string{"a", "b\nc", "d"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewCSVParser()
			actual, err := parser.ParseString(test.name)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expect, actual)
			}
		})
	}
}

func TestParser_ParseString_varyFields(t *testing.T) {
	parser := NewCSVParser()

	{
		fields, err := parser.ParseString("1,2")
		assert.NoError(t, err)
		assert.Equal(t, []string{"1", "2"}, fields)
	}
	{
		fields, err := parser.ParseString("1,2,3")
		assert.Equal(t, ErrFieldCount, err)
		assert.Equal(t, []string{"1", "2", "3"}, fields)
	}
	{
		fields, err := parser.ParseString("4,5")
		assert.NoError(t, err)
		assert.Equal(t, []string{"4", "5"}, fields)
	}
	{
		fields, err := parser.ParseString("6")
		assert.Equal(t, ErrFieldCount, err)
		assert.Equal(t, []string{"6"}, fields)
	}
}
