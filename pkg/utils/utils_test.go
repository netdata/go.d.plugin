package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

func TestToMap(t *testing.T) {
	s := data{
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
	expected := map[string]int64{
		"a": 1,
		"b": 2,
		"d": 4,
		"e": 5,
		"f": 6,
		"g": 7,
	}

	assert.EqualValuesf(t, expected, ToMap(s), "value test")
	assert.EqualValuesf(t, expected, ToMap(&s), "ptr test")
}

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

func TestDuration_ConvertTo(t *testing.T) {
	d := Duration{time.Second}

	assert.Equal(
		t,
		1000,
		d.ConvertTo(time.Millisecond),
	)

	d.Duration = time.Second * 3

	assert.Equal(
		t,
		1000*1000*3,
		d.ConvertTo(time.Microsecond),
	)
}
