package matcher

import (
	"errors"
	"strings"
)

const (
	// FmtString is a string match format.
	FmtString = "string"
	// FmtGlob is a glob match format.
	FmtGlob = "glob"
	// FmtRegExp is a regex match format.
	FmtRegExp = "regexp"
	// FmtSimplePattern is a simple pattern match format
	// https://docs.netdata.cloud/libnetdata/simple_pattern/
	FmtSimplePattern = "simplepattern"

	// Separator is a separator between match format and expression.
	Separator = ":"
)

const (
	SymString = '='
	SymGlob   = '*'
	SymRegExp = '~'
	SymNeg    = '!'
)

var (
	// ErrUnsupportedMatcherFormat error for unsupported matcher format
	ErrUnsupportedMatcherFormat = errors.New("unsupported matcher method")
	// ErrUnsupportedMatcherSyntax error for unsupported matcher syntax
	ErrUnsupportedMatcherSyntax = errors.New("unsupported matcher syntax")

	errNotShortSyntax = errors.New("not short syntax")
)

// Matcher is an interface that wraps MatchString method.
type Matcher interface {
	Match(b []byte) bool
	MatchString(string) bool
}

// New create a matcher
func New(format string, expr string) (Matcher, error) {
	switch format {
	case FmtString:
		return NewStringMatcher(expr), nil
	case FmtGlob:
		return NewGlobMatcher(expr)
	case FmtRegExp:
		return NewRegExpMatcher(expr)
	case FmtSimplePattern:
		return NewSimplePatternsMatcher(expr)
	default:
		return nil, ErrUnsupportedMatcherFormat
	}
}

func NewStringMatcher(expr string) Matcher {
	return stringFullMatcher(expr)
}

// Parse parses line and returns appropriate matcher based on match format.
//
// Short syntax
//   [ '!' ] <symbol> { ' ' } <expr>
//   = my_value
//   * *.example.com
//   ~ [0-9]+
//   != my_value
//
// Long syntax
//   <name>:<expr>
//   string:my_value
//   glob:*.example.com
//   regexp:[0-9]+
func Parse(line string) (Matcher, error) {
	matcher, err := parseShortFormat(line)
	if err == nil {
		return matcher, nil
	}
	if err != errNotShortSyntax {
		return nil, err
	}
	return parseLongSyntax(line)
}

func parseShortFormat(line string) (Matcher, error) {
	var format string
	var neg bool
	switch line {
	case "", "!":
		return nil, ErrUnsupportedMatcherFormat
	}
	if line[0] == SymNeg {
		neg = true
		line = line[1:]
	}
	switch line[0] {
	case SymString:
		format = FmtString
	case SymGlob:
		format = FmtGlob
	case SymRegExp:
		format = FmtRegExp
	default:
		return nil, errNotShortSyntax
	}
	expr := line[1:]
	for i, c := range expr {
		if !isSpace(c) {
			expr = expr[i:]
			break
		}
	}
	m, err := New(format, expr)
	if err != nil {
		return nil, err
	}
	if neg {
		m = Not(m)
	}
	return m, nil
}

func isSpace(c rune) bool {
	switch c {
	case ' ', '\t', '\f', '\v':
		return true
	default:
		return false
	}
}

func parseLongSyntax(line string) (Matcher, error) {
	parts := strings.SplitN(line, Separator, 2)
	if len(parts) != 2 {
		return nil, ErrUnsupportedMatcherSyntax
	}
	return New(parts[0], parts[1])
}
