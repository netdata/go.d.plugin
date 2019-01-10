package matcher

import (
	"errors"
	"strings"
)

type MatchFormat string

const (
	FmtString    MatchFormat = "="
	FmtGlob      MatchFormat = "*"
	FmtRegExp    MatchFormat = "~"
	FmtNegString MatchFormat = "!="
	FmtNegGlob   MatchFormat = "!*"
	FmtNegRegExp MatchFormat = "!~"
)

const separator = ":"

// Matcher is an interface that wraps Match method.
type Matcher interface {
	Match(string) bool
}

type NegMatcher struct{ Matcher }

func (m NegMatcher) Match(line string) bool { return !m.Matcher.Match(line) }

// CreateMatcher creates matcher.
func CreateMatcher(format MatchFormat, expr string) (m Matcher, err error) {
	switch format {
	case FmtString, FmtNegString:
		m = createStringMatcher(expr)
	case FmtRegExp, FmtNegRegExp:
		m, err = createRegExpMatcher(expr)
	case FmtGlob, FmtNegGlob:
		m, err = createGlobMatcher(expr)
	}

	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, errors.New("unsupported matcher method")

	}

	if format[0] == '!' {

		return &NegMatcher{m}, nil

	}

	return m, nil

}

func Parse(line string) (Matcher, error) {
	parts := strings.SplitN(line, separator, 2)
	if len(parts) != 2 {
		return nil, errors.New("unsupported matcher syntax")
	}
	return CreateMatcher(MatchFormat(parts[0]), parts[1])
}
