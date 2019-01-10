package matcher

import "strings"

func createStringMatcher(expr string) Matcher {
	switch {
	case len(expr) > 2 && strings.HasPrefix(expr, "^") && strings.HasSuffix(expr, "$"):
		return &StringFull{expr[1 : len(expr)-1]}
	case strings.HasPrefix(expr, "^"):
		return &StringPrefix{expr[1:]}
	case strings.HasSuffix(expr, "$"):
		return &StringSuffix{expr[:len(expr)-1]}
	default:
		return &StringPartial{expr}
	}
}

// StringFull implements Matcher, it uses "==" to match.
type StringFull struct{ Str string }

// Match matches.
func (m StringFull) Match(line string) bool { return m.Str == line }

// StringPartial implements Matcher, it uses strings.Contains to match.
type StringPartial struct{ Substr string }

// Match matches.
func (m StringPartial) Match(line string) bool { return strings.Contains(line, m.Substr) }

// StringPrefix implements Matcher, it uses strings.HasPrefix to match.
type StringPrefix struct{ Prefix string }

// Match matches.
func (m StringPrefix) Match(line string) bool { return strings.HasPrefix(line, m.Prefix) }

// StringSuffix implements Matcher, it uses strings.HasSuffix to match.
type StringSuffix struct{ Suffix string }

// Match matches.
func (m StringSuffix) Match(line string) bool { return strings.HasSuffix(line, m.Suffix) }
