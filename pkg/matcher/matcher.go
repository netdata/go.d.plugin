package matcher

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

// MatchFormat match format.
type MatchFormat string

const (
	// FmtString is a string match format.
	FmtString MatchFormat = "="
	// FmtGlob is a glob match format.
	FmtGlob MatchFormat = "*"
	// FmtRegExp is a regex[ match format.
	FmtRegExp MatchFormat = "~"
	// FmtNegString is a negative string match format.
	FmtNegString MatchFormat = "!="
	// FmtNegGlob is a negative glob match format.
	FmtNegGlob MatchFormat = "!*"
	// FmtNegRegExp is a negative regexp match format.
	FmtNegRegExp MatchFormat = "!~"
)

// Separator is a separator between match format and expression.
const Separator = ":"

// Matcher is an interface that wraps Match method.
type Matcher interface {
	Match(string) bool
}

// NegMatcher is a Matcher wrapper. It returns negative match.
type NegMatcher struct{ Matcher }

// Match matches
func (m NegMatcher) Match(line string) bool { return !m.Matcher.Match(line) }

// Create creates matcher based on match format.
func Create(format MatchFormat, expr string) (m Matcher, err error) {
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
		m = &NegMatcher{m}
	}

	return m, nil
}

// Parses parses line and returns appropriate matcher based on match format.
func Parse(line string) (Matcher, error) {
	parts := strings.SplitN(line, Separator, 2)
	if len(parts) != 2 {
		return nil, errors.New("unsupported matcher syntax")
	}
	return Create(MatchFormat(parts[0]), parts[1])
}

func createStringMatcher(expr string) Matcher {
	full := len(expr) > 2 && strings.HasPrefix(expr, "^") && strings.HasSuffix(expr, "$")
	prefix := strings.HasPrefix(expr, "^")
	suffix := strings.HasSuffix(expr, "$")

	switch {
	case full:
		return &StringFull{expr[1 : len(expr)-1]}
	case prefix:
		return &StringPrefix{expr[1:]}
	case suffix:
		return &StringSuffix{expr[:len(expr)-1]}
	default:
		return &StringPartial{expr}
	}
}

func createGlobMatcher(expr string) (Matcher, error) {
	if _, err := filepath.Match(expr, "QQ"); err != nil {
		return nil, err
	}
	return &GlobMatch{expr}, nil
}

func createRegExpMatcher(expr string) (Matcher, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &RegExpMatch{re}, nil
}
