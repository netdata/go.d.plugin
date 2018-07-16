package web_log

import (
	"errors"
	"github.com/l2isbad/yaml"
	"regexp"
)

type category struct {
	id   string
	name string
	re   *regexp.Regexp
}

type categories struct {
	prefix string
	list   []*category
}

func (c *categories) other() string {
	return c.prefix + "_other"
}

func (c *categories) add(n string, r *regexp.Regexp) {
	c.list = append(c.list, &category{c.prefix + "_" + n, n, r})
}

func (c *categories) active() bool {
	return c.list != nil
}

type rawCategory struct {
	name string
	re   string
}

type rawCategories []rawCategory

func (c *rawCategories) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var m yaml.MapSlice

	if err := unmarshal(&m); err != nil {
		return err
	}
	for k := range m {
		v, ok := m[k].Value.(string)
		if !ok {
			return errors.New("\"categories\" bad format")
		}
		*c = append(*c, rawCategory{m[k].Key.(string), v})
	}
	return nil
}
