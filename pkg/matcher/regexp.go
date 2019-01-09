package matcher

import "regexp"

type RegexpMatch struct{ Regexp *regexp.Regexp }

func (m RegexpMatch) Match(line string) bool { return m.Regexp.MatchString(line) }
