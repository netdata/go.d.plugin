package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/matcher"

	"github.com/prometheus/prometheus/pkg/labels"
)

type Matcher interface {
	Matches(labels labels.Labels) bool
}

const (
	FmtEqual     = "="
	FmtNegEqual  = "!="
	FmtRegexp    = "=~"
	FmtNegRegexp = "!~"
	FmtGlob      = "=*"
	FmtNegGlob   = "!*"
)

var (
	reLV = regexp.MustCompile(`^(?P<label_name>[a-zA-Z0-9_]+)(?P<format>=~|!~|=\*|!\*|=|!=)"(?P<expr>.+)"$`)
)

type labelMatcher struct {
	name string
	m    matcher.Matcher
}

func (m labelMatcher) Matches(lbs labels.Labels) bool {
	if m.name == labels.MetricName {
		return m.m.MatchString(lbs[0].Value)
	}
	if label, ok := lookupLabel(m.name, lbs[1:]); ok {
		return m.m.MatchString(label.Value)
	}
	return false
}

func Parse(expr string) (Matcher, error) {
	var matchers []Matcher
	lvs := strings.Split(unsugarExpr(expr), ",")

	for _, lv := range lvs {
		submatch := reLV.FindStringSubmatch(strings.TrimSpace(lv))
		if submatch == nil {
			return nil, fmt.Errorf("invalid expr syntax: '%s'", lv)
		}

		name, format, value := submatch[1], submatch[2], strings.Trim(submatch[3], "\"")

		var m matcher.Matcher
		var err error

		switch format {
		case FmtEqual, FmtNegEqual:
			m, err = matcher.NewStringMatcher(value, true, true)
		case FmtRegexp, FmtNegRegexp:
			m, err = matcher.NewRegExpMatcher(value)
		case FmtGlob, FmtNegGlob:
			if len(strings.Fields(value)) > 1 {
				m, err = matcher.NewSimplePatternsMatcher(value)
			} else {
				m, err = matcher.NewGlobMatcher(value)
			}
		default:
			err = fmt.Errorf("unknown format: %s", format)
		}
		if err != nil {
			return nil, err
		}

		lm := labelMatcher{
			name: name,
			m:    m,
		}

		if neg := strings.HasPrefix(format, "!"); neg {
			matchers = append(matchers, Not(lm))
		} else {
			matchers = append(matchers, lm)
		}
	}

	switch len(matchers) {
	case 0:
		return nil, nil
	case 1:
		return matchers[0], nil
	default:
		return And(matchers[0], matchers[1], matchers[2:]...), nil
	}
}

func unsugarExpr(expr string) string {
	expr = strings.TrimSpace(expr)
	switch idx := strings.IndexByte(expr, '{'); true {
	case idx == -1:
		expr = fmt.Sprintf(`__name__%s"%s"`, FmtGlob, expr)
	case idx > 0:
		expr = fmt.Sprintf(`__name__%s"%s",%s`, FmtGlob, expr[:idx], strings.Trim(expr[idx:], "{}"))
	default:
		expr = strings.Trim(expr, "{}")
	}
	return expr
}

func lookupLabel(name string, lbs labels.Labels) (labels.Label, bool) {
	for _, label := range lbs {
		if label.Name == name {
			return label, true
		}
	}
	return labels.Label{}, false
}
