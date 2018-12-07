package utils

import (
	"reflect"
)

const fieldTagName = "stm"
const prefixTagName = "prefix"

// ToMap converts struct to a map[string]int64 based on 'stm' tags
func ToMap(s interface{}) map[string]int64 {
	rv := make(map[string]int64)
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()
	toMap(v, t, rv, "")
	return rv
}

func toMap(v reflect.Value, t reflect.Type, rv map[string]int64, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if ft.Type.Kind() == reflect.Struct {
			nestPrefix := prefix + ft.Tag.Get(prefixTagName)
			toMap(v.Field(i), ft.Type, rv, nestPrefix)
		}
		if tag := ft.Tag.Get(fieldTagName); tag != "" {
			rv[prefix+tag] = v.Field(i).Int()
		}
	}
}
