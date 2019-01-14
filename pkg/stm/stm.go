package stm

import (
	"reflect"

	"github.com/netdata/go.d.plugin/logger"
)

const fieldTagName = "stm"

// ToMap converts struct to a map[string]int64 based on 'stm' tags
func ToMap(s ...interface{}) map[string]int64 {
	rv := map[string]int64{}
	for _, v := range s {
		value := reflect.Indirect(reflect.ValueOf(v))
		toMap(value, rv, "")
	}
	return rv
}

func toMap(value reflect.Value, rv map[string]int64, key string) {
	switch value.Kind() {
	case reflect.Ptr:
		convertPtr(value, rv, key)
	case reflect.Struct:
		convertStruct(value, rv, key)
	case reflect.Map:
		convertMap(value, rv, key)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		convertInteger(value, rv, key)
	case reflect.Interface:
		convertInterface(value, rv, key)
	default:
		logger.Panic("unsupported data type: ", value.Kind())
	}
}

func convertPtr(value reflect.Value, rv map[string]int64, key string) {
	if !value.IsNil() {
		toMap(value.Elem(), rv, key)
	}
}

func convertStruct(value reflect.Value, rv map[string]int64, key string) {
	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		tag, ok := ft.Tag.Lookup(fieldTagName)
		if !ok {
			continue
		}
		value := value.Field(i)
		toMap(value, rv, joinPrefix(key, tag))
	}
}

func convertMap(value reflect.Value, rv map[string]int64, key string) {
	for _, k := range value.MapKeys() {
		toMap(value.MapIndex(k), rv, joinPrefix(key, k.String()))
	}
}

func convertInteger(value reflect.Value, rv map[string]int64, key string) {
	intVal := value.Int()
	if _, ok := rv[key]; ok {
		logger.Panic("duplicate key: ", key)
	}
	rv[key] = intVal
}

func convertInterface(value reflect.Value, rv map[string]int64, key string) {
	fv := reflect.ValueOf(value.Interface())
	toMap(fv, rv, key)
}

func joinPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "_" + key
}
