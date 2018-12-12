package weblog

import (
	"fmt"

	"github.com/netdata/go.d.plugin/modules/weblog/matcher"
)

type rawCategory struct {
	Name  string
	Match string
}

type category struct {
	name string
	matcher.Matcher
}

func newCategory(raw rawCategory) (*category, error) {
	cat := &category{}

	if raw.Name == "" || raw.Match == "" {
		return nil, fmt.Errorf("category bad syntax : %s", raw)
	}

	m, err := matcher.New(raw.Match)

	if err != nil {
		return nil, err
	}
	cat.Matcher = m
	cat.name = raw.Name

	return cat, nil
}
