package matcher

import (
	"errors"
	"fmt"
	"regexp"
)

type (
	// Matcher is an interface that wraps MatchString method.
	Matcher interface {
		// Perform match against given []byte
		Match(b []byte) bool
		// Perform match against given string
		MatchString(string) bool
	}

	// Format matcher format
	Format string
)

const (
	// FmtString is a string match format.
	FmtString Format = "string"
	// FmtGlob is a glob match format.
	FmtGlob Format = "glob"
	// FmtRegExp is a regex match format.
	FmtRegExp Format = "regexp"
	// FmtSimplePattern is a simple pattern match format
	// https://docs.netdata.cloud/libnetdata/simple_pattern/
	FmtSimplePattern Format = "simple_patterns"

	// Separator is a separator between match format and expression.
	Separator = ":"
)

const (
	symString = "="
	symGlob   = "*"
	symRegExp = "~"
)

var (
	reShortSyntax = regexp.MustCompile(`(?s)^(!)?(.)\s*(.*)$`)
	reLongSyntax  = regexp.MustCompile(`(?s)^(!)?([^:]+):(.*)$`)

	errNotShortSyntax = errors.New("not short syntax")
)

// New create a matcher
func New(format Format, expr string) (Matcher, error) {
	switch format {
	case FmtString:
		return NewStringMatcher(expr, true, true)
	case FmtGlob:
		return NewGlobMatcher(expr)
	case FmtRegExp:
		return NewRegExpMatcher(expr)
	case FmtSimplePattern:
		return NewSimplePatternsMatcher(expr)
	default:
		return nil, fmt.Errorf("unsupported matcher format: '%s'", format)
	}
}

// Parse parses line and returns appropriate matcher based on match format.
func Parse(line string) (Matcher, error) {
	matcher, err := parseShortFormat(line)
	if err == nil {
		return matcher, nil
	}
	return parseLongSyntax(line)
}

func parseShortFormat(line string) (Matcher, error) {
	m := reShortSyntax.FindStringSubmatch(line)
	if m == nil {
		return nil, errNotShortSyntax
	}
	var format Format
	switch m[2] {
	case symString:
		format = FmtString
	case symGlob:
		format = FmtGlob
	case symRegExp:
		format = FmtRegExp
	default:
		return nil, fmt.Errorf("invalid short syntax: unknown symbol '%s'", m[2])
	}
	expr := m[3]
	matcher, err := New(format, expr)
	if err != nil {
		return nil, err
	}
	if m[1] != "" {
		matcher = Not(matcher)
	}
	return matcher, nil
}

func parseLongSyntax(line string) (Matcher, error) {
	m := reLongSyntax.FindStringSubmatch(line)
	if m == nil {
		return nil, fmt.Errorf("invalid syntax")
	}
	matcher, err := New(Format(m[2]), m[3])
	if err != nil {
		return nil, err
	}
	if m[1] != "" {
		matcher = Not(matcher)
	}
	return matcher, nil
}
