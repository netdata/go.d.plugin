package stm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/netdata/go-orchestrator/logger"
)

const fieldTagName = "stm"

type (
	Value interface {
		WriteTo(rv map[string]int64, key string, mul, div int)
	}
)

// ToMap converts struct to a map[string]int64 based on 'stm' tags
func ToMap(s ...interface{}) map[string]int64 {
	rv := map[string]int64{}
	for _, v := range s {
		value := reflect.Indirect(reflect.ValueOf(v))
		toMap(value, rv, "", 1, 1)
	}
	return rv
}

func toMap(value reflect.Value, rv map[string]int64, key string, mul, div int) {
	if value.CanInterface() {
		val, ok := value.Interface().(Value)
		if ok {
			val.WriteTo(rv, key, mul, div)
			return
		}
	}
	switch value.Kind() {
	case reflect.Ptr:
		convertPtr(value, rv, key, mul, div)
	case reflect.Struct:
		convertStruct(value, rv, key)
	case reflect.Map:
		convertMap(value, rv, key, mul, div)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		convertInteger(value, rv, key, mul, div)
	case reflect.Float32, reflect.Float64:
		convertFloat(value, rv, key, mul, div)
	case reflect.Interface:
		convertInterface(value, rv, key, mul, div)
	default:
		logger.Panic("unsupported data type: ", value.Kind())
	}
}

func convertPtr(value reflect.Value, rv map[string]int64, key string, mul, div int) {
	if !value.IsNil() {
		toMap(value.Elem(), rv, key, mul, div)
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
		prefix, mul, div := parseTag(tag)
		toMap(value, rv, joinPrefix(key, prefix), mul, div)
	}
}

func convertMap(value reflect.Value, rv map[string]int64, key string, mul, div int) {
	for _, k := range value.MapKeys() {
		toMap(value.MapIndex(k), rv, joinPrefix(key, k.String()), mul, div)
	}
}

func convertInteger(value reflect.Value, rv map[string]int64, key string, mul, div int) {
	intVal := value.Int()
	if _, ok := rv[key]; ok {
		logger.Panic("duplicate key: ", key)
	}
	rv[key] = intVal * int64(mul) / int64(div)
}

func convertFloat(value reflect.Value, rv map[string]int64, key string, mul, div int) {
	floatVal := value.Float()
	if _, ok := rv[key]; ok {
		logger.Panic("duplicate key: ", key)
	}
	rv[key] = int64(floatVal * float64(mul) / float64(div))
}

func convertInterface(value reflect.Value, rv map[string]int64, key string, mul, div int) {
	fv := reflect.ValueOf(value.Interface())
	toMap(fv, rv, key, mul, div)
}

func joinPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "_" + key
}

func parseTag(tag string) (prefix string, mul int, div int) {
	tokens := strings.Split(tag, ",")
	mul = 1
	div = 1
	var err error
	switch len(tokens) {
	case 3:
		div, err = strconv.Atoi(tokens[2])
		if err != nil {
			logger.Panic(err)
		}
		fallthrough
	case 2:
		mul, err = strconv.Atoi(tokens[1])
		if err != nil {
			logger.Panic(err)
		}
		fallthrough
	case 1:
		prefix = tokens[0]
	default:
		logger.Panic(fmt.Errorf("invalid tag format: %s", tag))
	}
	return
}
