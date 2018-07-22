package flags

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type flagSet interface {
	Name() string
	ShortName() string
	UserDefaultValue() interface{}
	Usage() string

	NeedArgument() bool
	Set(interface{})
	SetDefault()
}

type FlagContext interface {
	Parse()
	StringVar(*string, string, string, string, string)
	BoolVar(*bool, string, string, bool, string)
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
		flags: make([]flagSet, 0),
		args:  a,
	}
}

type flagContext struct {
	flags []flagSet
	args  []string
}

func (fc *flagContext) StringVar(v *string, name string, shortName string, value string, usage string) {
	f := &stringFlag{
		value:        v,
		name:         name,
		shortName:    shortName,
		defaultValue: value,
		usage:        usage,
	}
	fc.flags = append(fc.flags, f)
}

func (fc *flagContext) BoolVar(v *bool, name string, shortName string, value bool, usage string) {
	f := &boolFlag{
		value:        v,
		name:         name,
		shortName:    shortName,
		defaultValue: value,
		usage:        usage,
	}
	fc.flags = append(fc.flags, f)
}

func (fc flagContext) Parse() {
	if hasHelp(fc.args) {
		fmt.Println(help(fc.flags))
		os.Exit(2)
	}

	err := fc.parse()

	if err != nil {
		fmt.Println(err, "\n", help(fc.flags))
		os.Exit(2)
	}
}

func (fc flagContext) parse() error {
	for _, f := range fc.flags {
		f.Set(f.UserDefaultValue())
	}

	for i := 0; i < len(fc.args); i++ {

		if !isFlag(fc.args[i]) {
			continue
		}

		v := strings.TrimLeft(fc.args[i], "-")

		f, ok := fc.lookupFlags(v)
		if !ok {
			return fmt.Errorf("flag provided but not defined: %s", fc.args[i])
		}

		switch {
		default:
			f.Set(fc.args[i+1])
			i++
		case i == len(fc.args)-1, isFlag(fc.args[i+1]):
			if f.NeedArgument() {
				return fmt.Errorf("flag needs an argument: %s", fc.args[i])
			}
			f.SetDefault()
		}
	}
	return nil
}

func (fc flagContext) lookupFlags(v string) (flagSet, bool) {
	for _, f := range fc.flags {
		if v == f.Name() || v == f.ShortName() {
			return f, true
		}
	}
	return nil, false
}

// TODO: this is lame(negative integers...)
func isFlag(v string) bool {
	return strings.HasPrefix(v, "-")
}

func hasHelp(s []string) bool {
	for i := range s {
		v := strings.TrimLeft(s[i], "-")
		if v == "help" || v == "h" {
			return true
		}
	}
	return false
}

func help(s []flagSet) string {
	var msg = "Usage:\n"
	for _, f := range s {
		msg += fmt.Sprintf(
			"  --%s (-%s) %T\n\t%s (default %v)\n",
			f.Name(),
			f.ShortName(),
			f.UserDefaultValue(),
			f.Usage(),
			f.UserDefaultValue(),
		)
	}
	return msg
}
