package matcher

import "regexp"

// RegExpMatch implements Matcher, it uses regexp.MatchString to match.
type RegExpMatch struct{ RegExp *regexp.Regexp }

// Match matches.
func (m RegExpMatch) Match(line string) bool { return m.RegExp.MatchString(line) }
