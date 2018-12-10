package utils

import (
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

func TestDuration_UnmarshalYAML(t *testing.T) {
	var d Duration
	values := [][]byte{
		[]byte("100ms"),   // duration
		[]byte("3s300ms"), // duration
		[]byte("3"),       // int
		[]byte("3.3"),     // float
	}

	for _, v := range values {
		assert.NoError(t, yaml.Unmarshal(v, &d))
	}
}
