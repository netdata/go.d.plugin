package utils

import "reflect"

const fieldTagName = "stm"

// ToMap Convert struct to map[string]int64 based on stm tag
func ToMap(s interface{}) map[string]int64 {
	rv := make(map[string]int64)
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()
	toMap(v, t, rv)
	return rv
}

func toMap(v reflect.Value, t reflect.Type, rv map[string]int64) {
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if ft.Type.Kind() == reflect.Struct {
			toMap(v.Field(i), ft.Type, rv)
		}
		if tag := ft.Tag.Get(fieldTagName); tag != "" {
			rv[tag] = v.Field(i).Int()
		}
	}
}