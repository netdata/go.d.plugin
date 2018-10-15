package web_log

import (
	"gopkg.in/yaml.v2"

	"github.com/l2isbad/go.d.plugin/internal/modules/web_log/matcher"
)

type categories struct {
	items []category
	other string
}

func (c categories) exist() bool {
	return len(c.items) > 0
}

type category struct {
	id   string
	name string
	matcher.Matcher
}

func getCategories(ms yaml.MapSlice, prefix string) (categories, error) {
	cats := categories{
		other: prefix + "_other",
	}

	if len(ms) == 0 {
		return cats, nil
	}

	for _, v := range ms {
		r, ok := v.Value.(string)
		if !ok || r == "" {
			continue
		}

		m, err := matcher.New(r)
		if err != nil {
			return cats, err
		}

		cat := category{
			id:      prefix + "_" + v.Key.(string),
			name:    v.Key.(string),
			Matcher: m,
		}
		cats.items = append(cats.items, cat)
	}
	return cats, nil
}
