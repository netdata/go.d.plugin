package funcs

import (
	"reflect"
)

var fieldTagName = "stm"

// ToMap returns *map[string]int obtained by converting input *struct to map[string]int.
// Keys for the map function takes from "stm" tags. All tagged fields must be int type.
// Function panics if input value not a reference to a struct.
func ToMap(ptrStruct interface{}) map[string]int64 {
	rv := make(map[string]int64)
	t, v := reflect.TypeOf(ptrStruct).Elem(), reflect.ValueOf(ptrStruct).Elem()
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if tag := ft.Tag.Get(fieldTagName); tag != "" {
			rv[tag] = v.Field(i).Int()
		}
	}
	return rv
}
