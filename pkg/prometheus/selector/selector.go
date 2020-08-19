package selector

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/matcher"

	"github.com/prometheus/prometheus/pkg/labels"
)

type Selector interface {
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

type labelSelector struct {
	name string
	m    matcher.Matcher
}

func (s labelSelector) Matches(lbs labels.Labels) bool {
	if s.name == labels.MetricName {
		return s.m.MatchString(lbs[0].Value)
	}
	if label, ok := lookupLabel(s.name, lbs[1:]); ok {
		return s.m.MatchString(label.Value)
	}
	return false
}

func Parse(expr string) (Selector, error) {
	var srs []Selector
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

		sr := labelSelector{
			name: name,
			m:    m,
		}

		if neg := strings.HasPrefix(format, "!"); neg {
			srs = append(srs, Not(sr))
		} else {
			srs = append(srs, sr)
		}
	}

	switch len(srs) {
	case 0:
		return nil, nil
	case 1:
		return srs[0], nil
	default:
		return And(srs[0], srs[1], srs[2:]...), nil
	}
}

func unsugarExpr(expr string) string {
	expr = strings.TrimSpace(expr)
	var metricName, matchers string

	if idx := strings.IndexByte(expr, '{'); idx == -1 {
		metricName = expr
	} else if idx == 0 {
		matchers = strings.Trim(expr, "{}")
	} else {
		metricName, matchers = expr[:idx], strings.Trim(expr[idx:], "{}")
	}
	metricName, matchers = strings.TrimSpace(metricName), strings.TrimSpace(matchers)

	if metricName == "" {
		return matchers
	}

	metricNameOp := FmtGlob
	if len(strings.Fields(metricName)) == 1 && metricName[0] == '!' {
		metricName = metricName[1:]
		metricNameOp = FmtNegGlob
	}

	if matchers == "" {
		return fmt.Sprintf(`__name__%s"%s"`, metricNameOp, metricName)
	}
	return fmt.Sprintf(`__name__%s"%s",%s`, metricNameOp, metricName, matchers)
}

func lookupLabel(name string, lbs labels.Labels) (labels.Label, bool) {
	for _, label := range lbs {
		if label.Name == name {
			return label, true
		}
	}
	return labels.Label{}, false
}
