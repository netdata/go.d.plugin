package notsimplepattern

import (
	"path/filepath"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/matcher"
)

// Pattern pattern.
type Pattern struct {
	Exclude bool
	matcher.GlobMatch
}

// New creates new not simple pattern.
func New() *Patterns {
	return &Patterns{cache: make(map[string]bool)}
}

// Patterns patterns.
type Patterns struct {
	UseCache bool
	Patterns []Pattern

	cache map[string]bool
}

// Add adds pattern to the collections. The only possible returned error is ErrBadPattern.
func (ps *Patterns) Add(pattern string) error {
	if _, err := filepath.Match(pattern, "QQ"); err != nil {
		return err
	}

	p := Pattern{}

	if strings.HasPrefix(pattern, "!") {
		p.Exclude = true
		p.Pattern = pattern[1:]
	} else {
		p.Pattern = pattern
	}

	ps.Patterns = append(ps.Patterns, p)

	return nil
}

// Match matches.
func (ps Patterns) Match(line string) bool {
	for _, p := range ps.Patterns {
		if p.Match(line) {
			return !p.Exclude
		}
	}
	return false
}

func Create(expr string) (*Patterns, error) {
	ps := New()

	for _, pattern := range strings.Fields(expr) {
		if err := ps.Add(pattern); err != nil {
			return nil, err
		}
	}

	return ps, nil
}
