package web_log

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"strings"
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
	c.list = append(c.list, &category{c.prefix + "_" + n,n, r})
}

func (c *categories) active() bool {
	return c.list != nil
}

type rawCategory struct {
	name string
	re   string
}

type rawCategories []rawCategory

func (c *rawCategories) UnmarshalTOML(input []byte) error {
	s := bufio.NewScanner(bytes.NewBuffer(input))
	s.Scan()
	for s.Scan() {
		val := strings.SplitN(s.Text(), "=", 2)
		if len(val) != 2 {
			return errors.New("bad format")
		}
		n, r := strings.TrimSpace(val[0]), strings.Trim(val[1], "'\" ")
		*c = append(*c, rawCategory{n, r})
	}
	return nil
}
