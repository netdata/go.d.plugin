package flags

import (
	"os"
	"strconv"
	"strings"
)

type flagSet interface {
	getUsage() string
	getName() string
	getShortName() string
	getDefaultValue() interface{}
	set(interface{})
}

type FlagContext interface {
	Parse()
	StringVar(*string, string, string, string, string)
	BoolVar(*bool, string, string, bool, string)
	// ShowUsage() string
}

// TODO: usage
func New() FlagContext {
	var a []string
	last := len(os.Args) - 1

	if _, err := strconv.Atoi(os.Args[last]); err != nil {
		a = os.Args[1:]
	} else {
		a = os.Args[1:last]
	}

	return &flagContext{
		flags: make(map[string]flagSet),
		args:  a,
	}
}

type flagContext struct {
	flags map[string]flagSet
	args  []string
}

func (fc *flagContext) StringVar(v *string, name string, shortName string, value string, usage string) {
	fc.flags[name] = &stringFlag{
		value:        v,
		name:         name,
		shortName:    shortName,
		defaultValue: value,
		usage:        usage,
	}
}

func (fc *flagContext) BoolVar(v *bool, name string, shortName string, value bool, usage string) {
	fc.flags[name] = &boolFlag{
		value:        v,
		name:         name,
		shortName:    shortName,
		defaultValue: value,
		usage:        usage,
	}
}

func (fc *flagContext) Parse() {

	for _, f := range fc.flags {
		idx := fc.inArgsIndex(f)

		if idx == -1 {
			f.set(f.getDefaultValue())
			continue
		}

		if len(fc.args) > idx+1 && isValue(fc.args[idx+1]) {
			f.set(fc.args[idx+1])
			continue
		}

		switch f.(type) {
		default:
			f.set(f.getDefaultValue())
		case *boolFlag:
			f.set(true)
		}
	}

}

func (fc *flagContext) inArgsIndex(f flagSet) int {
	for idx, v := range fc.args {
		if !isFlag(v) {
			continue
		}
		v = strings.TrimLeft(v, "-")
		if v == f.getName() || v == f.getShortName() {
			return idx
		}
	}
	return -1
}

func isFlag(v string) bool {
	return strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--")
}

func isValue(v string) bool {
	return !isFlag(v)
}
