package matcher

import (
	"path/filepath"
	"strings"
)

type ShellMatch struct{ Pattern string }

func (m ShellMatch) Match(line string) bool {
	ok, _ := filepath.Match(m.Pattern, line)
	return ok
}

type SimplePattern struct {
	Negative bool
	ShellMatch
}

type SimplePatterns []SimplePattern

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

	*s = append(*s, sp)

	return nil
}

func (s SimplePatterns) Match(line string) bool {
	for _, p := range s {
		if p.Match(line) {
			return p.Negative
		}
	}
	return false
}

func CreateSimplePatterns(line string) (*SimplePatterns, error) {
	sps := make(SimplePatterns, 0)

	for _, pattern := range strings.Fields(line) {

		if err := sps.Add(pattern); err != nil {
			return nil, err
		}
	}

	return &sps, nil
}

func checkShellPattern(pattern string) error {
	_, err := filepath.Match(pattern, "")
	return err
}
