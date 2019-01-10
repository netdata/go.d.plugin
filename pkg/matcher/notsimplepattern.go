package matcher

import (
	"strings"
)

type NotSimplePattern struct {
	Exclude bool
	GlobMatch
}

// NotSimplePatterns implements Matcher, it is an ordered collection of patterns.
type NotSimplePatterns []NotSimplePattern

// Add adds pattern to the collections. The only possible returned error is ErrBadPattern.
func (nsp *NotSimplePatterns) Add(pat string) error {
	//if err := checkGlobPatterns(pattern); err != nil {
	//	return err
	//}

	p := NotSimplePattern{}

	if strings.HasPrefix(pat, "!") {
		p.Exclude = true
		p.Pattern = pat[1:]
	} else {
		p.Pattern = pat
	}

	*nsp = append(*nsp, p)

	return nil
}

// Match matches.
func (nsp NotSimplePatterns) Match(line string) bool {
	for _, p := range nsp {
		if p.Match(line) {
			return !p.Exclude
		}
	}
	return false
}

func CreateSimplePatterns(expr string) (*NotSimplePatterns, error) {
	nsp := make(NotSimplePatterns, 0)

	for _, pattern := range strings.Fields(expr) {
		if err := nsp.Add(pattern); err != nil {
			return nil, err
		}
	}
	return &nsp, nil
}
