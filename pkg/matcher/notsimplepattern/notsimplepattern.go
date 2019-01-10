package notsimplepattern

import (
	"path/filepath"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type pattern struct {
	exclude bool
	matcher.GlobMatch
}

// New creates new not simple pattern.
func New() *Patterns {
	return &Patterns{Cache: make(map[string]bool)}
}

// Patterns patterns.
type Patterns struct {
	UseCache bool
	Cache    map[string]bool

	patterns []pattern
}

func (ps *Patterns) add(pat string) error {
	if _, err := filepath.Match(pat, "QQ"); err != nil {
		return err
	}

	p := pattern{}

	if strings.HasPrefix(pat, "!") {
		p.exclude = true
		p.Pattern = pat[1:]
	} else {
		p.Pattern = pat
	}

	ps.patterns = append(ps.patterns, p)

	return nil
}

// Match matches.
func (ps Patterns) Match(line string) bool {
	if !ps.UseCache {
		return ps.match(line)
	}

	if v, ok := ps.Cache[line]; ok {
		return v
	}

	matched := ps.match(line)
	ps.Cache[line] = matched

	return matched
}

func (ps Patterns) match(line string) bool {
	for _, p := range ps.patterns {
		if p.Match(line) {
			return !p.exclude
		}
	}
	return false
}

// Create creates new not simple patterns. It returns error in case one of patterns has bad syntax.
func Create(expr string) (*Patterns, error) {
	ps := New()

	for _, pattern := range strings.Fields(expr) {
		if err := ps.add(pattern); err != nil {
			return nil, err
		}
	}

	return ps, nil
}
