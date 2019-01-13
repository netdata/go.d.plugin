package matcher

import (
	"strings"
)

type (
	term struct {
		Matcher
		exclude bool
	}

	// Patterns patterns.
	Patterns []term
)

// NewSimplePatternsMatcher creates new simple patterns. It returns error in case one of patterns has bad syntax.
func NewSimplePatternsMatcher(expr string) (Matcher, error) {
	var ps Patterns

	for _, pattern := range strings.Fields(expr) {
		if err := ps.add(pattern); err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func (ps *Patterns) add(pat string) error {
	if pat == "" {
		return nil
	}
	p := term{}
	if pat[0] == '!' {
		p.exclude = true
		pat = pat[1:]
	}
	m, err := NewGlobMatcher(pat)
	if err != nil {
		return err
	}

	p.Matcher = m
	*ps = append(*ps, p)

	return nil
}

func (ps Patterns) Match(b []byte) bool {
	return ps.MatchString(string(b))
}

// MatchString matches.
func (ps Patterns) MatchString(line string) bool {
	for _, p := range ps {
		matched := p.MatchString(line)
		if !matched && p.exclude {
			return false
		}
		if matched && !p.exclude {
			return true
		}
	}
	return false
}
