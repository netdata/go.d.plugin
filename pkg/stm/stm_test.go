package stm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMap_empty(t *testing.T) {
	s := struct{}{}

	expected := map[string]int64{}

	assert.EqualValuesf(t, expected, ToMap(s), "value test")
	assert.EqualValuesf(t, expected, ToMap(&s), "ptr test")
}

func TestToMap_int(t *testing.T) {
	s := struct {
		I   int   `stm:"int"`
		I8  int8  `stm:"int8"`
		I16 int16 `stm:"int16"`
		I32 int32 `stm:"int32"`
		I64 int64 `stm:"int64"`
	}{
		I: 1, I8: 2, I16: 3, I32: 4, I64: 5,
	}

	expected := map[string]int64{
		"int": 1, "int8": 2, "int16": 3, "int32": 4, "int64": 5,
	}

	assert.EqualValuesf(t, expected, ToMap(s), "value test")
	assert.EqualValuesf(t, expected, ToMap(&s), "ptr test")
}

func TestToMap_struct(t *testing.T) {
	type pair struct {
		Left  int `stm:"left"`
		Right int `stm:"right"`
	}
	s := struct {
		I      int  `stm:"int"`
		Pempty pair `stm:""`
		Ps     pair `stm:"s"`
		Notag  int
	}{
		I:      1,
		Pempty: pair{2, 3},
		Ps:     pair{4, 5},
		Notag:  6,
	}

	expected := map[string]int64{
		"int":  1,
		"left": 2, "right": 3,
		"s_left": 4, "s_right": 5,
	}

	assert.EqualValuesf(t, expected, ToMap(s), "value test")
	assert.EqualValuesf(t, expected, ToMap(&s), "ptr test")
}

func TestToMap_tree(t *testing.T) {
	type node struct {
		Value int   `stm:"v"`
		Left  *node `stm:"left"`
		Right *node `stm:"right"`
	}
	s := node{1,
		&node{2, nil, nil},
		&node{3,
			&node{4, nil, nil},
			nil,
		},
	}
	expected := map[string]int64{
		"v":            1,
		"left_v":       2,
		"right_v":      3,
		"right_left_v": 4,
	}

	assert.EqualValuesf(t, expected, ToMap(s), "value test")
	assert.EqualValuesf(t, expected, ToMap(&s), "ptr test")
}

func TestToMap_map(t *testing.T) {
	s := struct {
		I int              `stm:"int"`
		M map[string]int64 `stm:""`
	}{
		I: 1,
		M: map[string]int64{
			"a": 2,
			"b": 3,
		},
	}

	expected := map[string]int64{
		"int": 1,
		"a":   2,
		"b":   3,
	}

	assert.EqualValuesf(t, expected, ToMap(s), "value test")
	assert.EqualValuesf(t, expected, ToMap(&s), "ptr test")
}

func TestToMap_nestMap(t *testing.T) {
	s := struct {
		I int                    `stm:"int"`
		M map[string]interface{} `stm:""`
	}{
		I: 1,
		M: map[string]interface{}{
			"a": 2,
			"b": 3,
			"m": map[string]interface{}{
				"c": 4,
			},
		},
	}

	expected := map[string]int64{
		"int": 1,
		"a":   2,
		"b":   3,
		"m_c": 4,
	}

	assert.EqualValuesf(t, expected, ToMap(s), "value test")
	assert.EqualValuesf(t, expected, ToMap(&s), "ptr test")
}

func TestToMap_ptr(t *testing.T) {
	two := 2
	s := struct {
		I   int  `stm:"int"`
		Ptr *int `stm:"ptr"`
		Nil *int `stm:"nil"`
	}{
		I:   1,
		Ptr: &two,
		Nil: nil,
	}

	expected := map[string]int64{
		"int": 1,
		"ptr": 2,
	}

	assert.EqualValuesf(t, expected, ToMap(s), "value test")
	assert.EqualValuesf(t, expected, ToMap(&s), "ptr test")
}

func TestToMap_invalidType(t *testing.T) {
	s := struct {
		Str string `stm:"int"`
	}{
		Str: "abc",
	}

	assert.Panics(t, func() {
		ToMap(s)
	}, "value test")
	assert.Panics(t, func() {
		ToMap(&s)
	}, "ptr test")
}

func TestToMap_duplicateKey(t *testing.T) {
	s := struct {
		Key int            `stm:"key"`
		M   map[string]int `stm:""`
	}{
		Key: 1,
		M: map[string]int{
			"key": 2,
		},
	}

	assert.Panics(t, func() {
		ToMap(s)
	}, "value test")
	assert.Panics(t, func() {
		ToMap(&s)
	}, "ptr test")
}

func TestToMap_Variadic(t *testing.T) {
	s1 := struct {
		Key1 int `stm:"key1"`
	}{
		Key1: 1,
	}
	s2 := struct {
		Key2 int `stm:"key2"`
	}{
		Key2: 2,
	}
	s3 := struct {
		Key3 int `stm:"key3"`
	}{
		Key3: 3,
	}

	assert.Equal(
		t,
		map[string]int64{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		},
		ToMap(s1, s2, s3),
	)
}
