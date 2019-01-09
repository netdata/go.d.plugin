package matcher

import "regexp"

// RegexpMatch implements Matcher, it uses regexp.MatchString to match.
type RegexpMatch struct{ Regexp *regexp.Regexp }

// Match matches.
func (m RegexpMatch) Match(line string) bool { return m.Regexp.MatchString(line) }
