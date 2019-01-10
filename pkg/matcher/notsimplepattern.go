package matcher

import (
	"strings"
)

type NotSimplePattern struct {
	Exclude bool
	GlobMatch
}

// NotSimplePatterns implements Matcher, it is an ordered collection of NotSimplePatterns.
type NotSimplePatterns struct {
	UseCache bool
	Patterns []NotSimplePattern

	cache map[string]bool
}

// Add adds pattern to the collections. The only possible returned error is ErrBadPattern.
func (s *NotSimplePatterns) Add(pattern string) error {
	if err := checkGlobPatterns(pattern); err != nil {
		return err
	}

	gp := NotSimplePattern{}

	if strings.HasPrefix(pattern, "!") {
		gp.Exclude = true
		gp.Pattern = pattern[1:]
	} else {
		gp.Pattern = pattern
	}

	s.Patterns = append(s.Patterns, gp)

	return nil
}

// Match matches.
func (s NotSimplePatterns) Match(line string) bool {
	if !s.UseCache {
		return s.match(line)
	}

	if v, ok := s.cache[line]; ok {
		return v
	}

	matched := s.match(line)
	s.cache[line] = matched

	return matched
}

func (s NotSimplePatterns) match(line string) bool {
	for _, p := range s.Patterns {
		if p.Match(line) {
			if p.Exclude {
				return false
			}
			return true
		}
	}
	return false
}

func CreateNotSimplePatterns(line string) (*NotSimplePatterns, error) {
	sps := &NotSimplePatterns{UseCache: true, cache: make(map[string]bool)}
	for _, pattern := range strings.Fields(line) {

		if err := sps.Add(pattern); err != nil {
			return nil, err
		}
	}
	return sps, nil
}
