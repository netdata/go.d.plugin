package flags

type stringFlag struct {
	value        *string
	name         string
	shortName    string
	defaultValue string
	usage        string
}

func (f *stringFlag) Set(i interface{}) {
	if v, ok := i.(string); ok {
		*f.value = v
	}
}

func (f stringFlag) Usage() string {
	return f.usage
}

func (f stringFlag) Name() string {
	return f.name
}

func (f stringFlag) ShortName() string {
	return f.shortName
}

func (f stringFlag) UserDefaultValue() interface{} {
	return f.defaultValue
}

func (f stringFlag) SetDefault() {
}

func (f stringFlag) NeedArgument() bool {
	return true
}
