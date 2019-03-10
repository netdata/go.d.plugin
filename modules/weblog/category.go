package weblog

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type RawCategory struct {
	Name  string `yaml:"name"`
	Match string `yaml:"match"`
}

type Category struct {
	name    string
	Matcher matcher.Matcher
}

func (r RawCategory) String() string {
	return fmt.Sprintf(`{"name": %q, "match": %q}`, r.Name, r.Match)
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
