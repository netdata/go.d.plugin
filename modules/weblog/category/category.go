package category

import (
	"fmt"

	"github.com/netdata/go.d.plugin/modules/weblog/matcher"
)

type Category interface {
	Name() string
	matcher.Matcher
}

type Raw struct {
	Name  string
	Match string
}

type category struct {
	name string
	matcher.Matcher
}

func (c category) Name() string {
	return c.name
}

func New(raw Raw) (Category, error) {
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
