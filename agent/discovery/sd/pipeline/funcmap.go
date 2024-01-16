// SPDX-License-Identifier: GPL-3.0-or-later

package pipeline

import (
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/netdata/go.d.plugin/pkg/matcher"
)

func newFuncMap() template.FuncMap {
	custom := map[string]interface{}{
		"match": funcMatchAny,
		"glob": func(value, pattern string, patterns ...string) bool {
			return funcMatchAny("glob", value, pattern, patterns...)
		},
	}

	fm := sprig.HermeticTxtFuncMap()

	for name, fn := range custom {
		fm[name] = fn
	}

	return fm
}

func funcMatchAny(typ, value, pattern string, patterns ...string) bool {
	switch len(patterns) {
	case 0:
		return funcMatch(typ, value, pattern)
	default:
		return funcMatch(typ, value, pattern) || funcMatchAny(typ, value, patterns[0], patterns[1:]...)
	}
}

func funcMatch(typ string, value, pattern string) bool {
	switch typ {
	case "glob", "":
		m, err := matcher.NewGlobMatcher(pattern)
		return err == nil && m.MatchString(value)
	case "sp":
		m, err := matcher.NewSimplePatternsMatcher(pattern)
		return err == nil && m.MatchString(value)
	case "re":
		ok, err := regexp.MatchString(pattern, value)
		return err == nil && ok
	case "dstar":
		ok, err := doublestar.Match(pattern, value)
		return err == nil && ok
	default:
		return false
	}
}
