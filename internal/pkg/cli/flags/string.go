package flags

type stringFlag struct {
	value        *string
	name         string
	shortName    string
	defaultValue string
	usage        string
}

func (f *stringFlag) set(i interface{}) {
	if v, ok := i.(string); ok {
		*f.value = v
	}
}

func (f *stringFlag) getUsage() string {
	return f.usage
}

func (f *stringFlag) getName() string {
	return f.name
}

func (f *stringFlag) getShortName() string {
	return f.shortName
}

func (f *stringFlag) getDefaultValue() interface{} {
	return f.defaultValue
}
