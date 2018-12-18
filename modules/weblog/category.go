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

func newCategory(raw rawCategory) (*category, error) {
	cat := &category{}

	if raw.Name == "" || raw.Match == "" {
		return nil, fmt.Errorf("category bad syntax : %s", raw)
	}

	m, err := newMatcher(raw.Match)

	if err != nil {
		return nil, err
	}
	cat.matcher = m

	return cat, nil
}
