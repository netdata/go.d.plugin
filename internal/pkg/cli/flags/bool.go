package flags

import "strconv"

type boolFlag struct {
	value        *bool
	name         string
	shortName    string
	defaultValue bool
	usage        string
}

func (f *boolFlag) Set(i interface{}) {
	switch v := i.(type) {
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			*f.value = b
		}
	case bool:
		*f.value = v
	}
}

func (f boolFlag) Usage() string {
	return f.usage
}

func (f boolFlag) Name() string {
	return f.name
}

func (f boolFlag) ShortName() string {
	return f.shortName
}

func (f boolFlag) UserDefaultValue() interface{} {
	return f.defaultValue
}

func (f *boolFlag) SetDefault() {
	*f.value = true
}

func (f *boolFlag) NeedArgument() bool {
	return false
}
