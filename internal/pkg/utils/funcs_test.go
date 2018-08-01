package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	Foo int64 `stm:"foo"`
	Bar int64 `stm:"bar"`
	Baz int64
}

func TestStrToMap_ptr(t *testing.T) {
	s := &data{1, 2, 3}
	m := StrToMap(s)

	assert.EqualValues(t, map[string]int64{
		"foo": 1,
		"bar": 2,
	}, m)
}

func TestStrToMap_value(t *testing.T) {
	s := data{1, 2, 3}
	m := StrToMap(s)

	assert.EqualValues(t, map[string]int64{
		"foo": 1,
		"bar": 2,
	}, m)
}
