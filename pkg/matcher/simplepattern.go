package matcher

import (
	"path/filepath"
	"strings"
)

// ShellMatch implements Matcher, it uses filepath.Match to match.
type ShellMatch struct{ Pattern string }

// Match matches.
func (m ShellMatch) Match(line string) bool {
	ok, _ := filepath.Match(m.Pattern, line)
	return ok
}

// SimplePattern implements Matcher, it uses filepath.Match to match.
type SimplePattern struct {
	Negative bool
	ShellMatch
}

// SimplePatterns implements Matcher, it is an ordered collection of SimplePatterns.
type SimplePatterns struct {
	UseCache bool
	Patterns []SimplePattern

	cache map[string]bool
}

// Add adds pattern to the collections. The only possible returned error is ErrBadPattern.
func (s *SimplePatterns) Add(pattern string) error {
	if err := checkShellPattern(pattern); err != nil {
		return err
	}

	sp := SimplePattern{}

	if strings.HasPrefix(pattern, "!") {
		sp.Negative = true
		sp.Pattern = pattern[1:]
	} else {
		sp.Pattern = pattern
	}

	s.Patterns = append(s.Patterns, sp)

	return nil
}

// Match matches.
func (s SimplePatterns) Match(line string) bool {
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

func (s SimplePatterns) match(line string) bool {
	for _, p := range s.Patterns {
		if p.Match(line) {
			if p.Negative {
				return false
			}
			return true
		}
	}
	return false
}

// CreateSimplePatterns creates cached SimplePatterns.
func CreateSimplePatterns(line string) (*SimplePatterns, error) {
	sps := &SimplePatterns{UseCache: true, cache: make(map[string]bool)}

	for _, pattern := range strings.Fields(line) {

		if err := sps.Add(pattern); err != nil {
			return nil, err
		}
	}

	return sps, nil
}

func checkShellPattern(pattern string) error {
	_, err := filepath.Match(pattern, "QQ")
	return err
}
