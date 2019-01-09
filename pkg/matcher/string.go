package matcher

import "strings"

// StringContains implements Matcher, it uses strings.Contains to match.
type StringContains struct{ Substr string }

// Match matches.
func (m StringContains) Match(line string) bool { return strings.Contains(line, m.Substr) }

// StringPrefix implements Matcher, it uses strings.HasPrefix to match.
type StringPrefix struct{ Prefix string }

// Match matches.
func (m StringPrefix) Match(line string) bool { return strings.HasPrefix(line, m.Prefix) }

// StringSuffix implements Matcher, it uses strings.HasSuffix to match.
type StringSuffix struct{ Suffix string }

// Match matches.
func (m StringSuffix) Match(line string) bool { return strings.HasSuffix(line, m.Suffix) }
