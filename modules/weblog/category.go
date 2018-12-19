package weblog

import (
	"fmt"
)

type rawcategory struct {
	Name  string
	Match string
}

type category struct {
	name string
	matcher
}

func (r rawcategory) String() string {
	return fmt.Sprintf("{name: %s, match: %s}", r.Name, r.Match)
}

func newCategory(raw rawcategory) (*category, error) {
	if raw.Name == "" || raw.Match == "" {
		return nil, fmt.Errorf("category bad syntax")
	}

	m, err := newMatcher(raw.Match)

	if err != nil {
		return nil, err
	}

	return &category{name: raw.Name, matcher: m}, nil
}
