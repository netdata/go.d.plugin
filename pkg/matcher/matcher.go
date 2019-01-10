package matcher

import (
	"errors"
	"path/filepath"
	"regexp"
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

type CachedMatcher struct {
	Cache map[string]bool
	Matcher
}

func (m CachedMatcher) Match(line string) bool {
	if v, ok := m.Cache[line]; ok {
		return v
	}

	matched := m.Matcher.Match(line)
	m.Cache[line] = matched

	return matched
}

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
		m = &NegMatcher{m}
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
	if err := checkGlobPatterns(expr); err != nil {
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

func checkGlobPatterns(pattern string) error {
	_, err := filepath.Match(pattern, "QQ")
	return err
}
