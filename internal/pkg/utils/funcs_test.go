package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	A int64 `stm:"a"`
	B int64 `stm:"b"`
	C int64
	X inner1
	inner2
}

type inner1 struct {
	D int64 `stm:"d"`
	E int64 `stm:"e"`
}

type inner2 struct {
	F int64 `stm:"f"`
	G int64 `stm:"g"`
}

func TestStrToMap_ptr(t *testing.T) {
	s := &data{
		A: 1,
		B: 2,
		C: 3,
		X: inner1{
			D: 4,
			E: 5,
		},
		inner2: inner2{
			F: 6,
			G: 7,
		},
	}
	m := ToMap(s)

	assert.EqualValues(t, map[string]int64{
		"a": 1,
		"b": 2,
		"d": 4,
		"e": 5,
		"f": 6,
		"g": 7,
	}, m)
}

func TestStrToMap_value(t *testing.T) {
	s := &data{
		A: 1,
		B: 2,
		C: 3,
		X: inner1{
			D: 4,
			E: 5,
		},
		inner2: inner2{
			F: 6,
			G: 7,
		},
	}
	m := ToMap(s)

	assert.EqualValues(t, map[string]int64{
		"a": 1,
		"b": 2,
		"d": 4,
		"e": 5,
		"f": 6,
		"g": 7,
	}, m)
}
