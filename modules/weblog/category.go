package weblog

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type RawCategory struct {
	Name  string
	Match string
}

type Category struct {
	name string
	matcher.Matcher
}

func (r RawCategory) String() string {
	return fmt.Sprintf("{name: %s, match: %s}", r.Name, r.Match)
}

func NewCategory(raw RawCategory) (*Category, error) {
	if raw.Name == "" || raw.Match == "" {
		return nil, fmt.Errorf("category bad syntax")
	}

	m, err := matcher.Parse(raw.Match)

	if err != nil {
		return nil, err
	}

	return &Category{name: raw.Name, Matcher: m}, nil
}
