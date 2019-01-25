package weblog

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type rawcategory struct {
	Name  string
	Match string
}

type category struct {
	name string
	matcher.Matcher
}

func (r rawcategory) String() string {
	return fmt.Sprintf("{name: %s, match: %s}", r.Name, r.Match)
}

func newCategory(raw rawcategory) (*category, error) {
	if raw.Name == "" || raw.Match == "" {
		return nil, fmt.Errorf("category bad syntax")
	}

	m, err := matcher.Parse(raw.Match)

	if err != nil {
		return nil, err
	}

	return &category{name: raw.Name, Matcher: m}, nil
}
