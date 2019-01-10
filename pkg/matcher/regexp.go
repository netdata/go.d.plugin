package matcher

import "regexp"

func createRegExpMatcher(expr string) (Matcher, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &RegExpMatch{re}, nil
}

// RegExpMatch implements Matcher, it uses regexp.MatchString to match.
type RegExpMatch struct{ RegExp *regexp.Regexp }

// Match matches.
func (m RegExpMatch) Match(line string) bool { return m.RegExp.MatchString(line) }
