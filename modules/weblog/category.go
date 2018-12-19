package weblog

import (
	"fmt"
)

type rawCategory struct {
	Name  string
	Match string
}

type category struct {
	name string
	matcher
}

func (r rawCategory) String() string {
	return fmt.Sprintf("{name: %s, match: %s}", r.Name, r.Match)
}

func newCategory(raw rawCategory) (*category, error) {
	if raw.Name == "" || raw.Match == "" {
		return nil, fmt.Errorf("category bad syntax")
	}

	m, err := newMatcher(raw.Match)

	if err != nil {
		return nil, err
	}

	return &category{name: raw.Name, matcher: m}, nil
}
