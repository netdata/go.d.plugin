package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	{
		args := []string{"go.d.plugin"}
		expected := &Option{Module: "all", UpdateEvery: 1}
		actual, err := Parse(args)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	{
		args := []string{"go.d.plugin", "2"}
		expected := &Option{Module: "all", UpdateEvery: 2}
		actual, err := Parse(args)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	{
		args := []string{"go.d.plugin", "2", "-d"}
		expected := &Option{Debug: true, Module: "all", UpdateEvery: 2}
		actual, err := Parse(args)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	{
		args := []string{"go.d.plugin", "2", "-m", "foo"}
		expected := &Option{Module: "foo", UpdateEvery: 2}
		actual, err := Parse(args)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	{
		args := []string{"go.d.plugin", "NaN", "-m", "foo"}
		_, err := Parse(args)
		assert.Error(t, err)
	}

	{
		args := []string{"go.d.plugin", "-m"}
		_, err := Parse(args)
		assert.Error(t, err)
	}
}
