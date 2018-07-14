package flags

import "strconv"

type boolFlag struct {
	value        *bool
	name         string
	shortName    string
	defaultValue bool
	usage        string
}

func (f *boolFlag) set(i interface{}) {
	switch v := i.(type) {
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			*f.value = b
		}
	case bool:
		*f.value = v
	}
}

func (f *boolFlag) getUsage() string {
	return f.usage
}

func (f *boolFlag) getName() string {
	return f.name
}

func (f *boolFlag) getShortName() string {
	return f.shortName
}

func (f *boolFlag) getDefaultValue() interface{} {
	return f.defaultValue
}
