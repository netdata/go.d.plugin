package parser

import "regexp"

func NewRegexpParser(regexp *regexp.Regexp) Parser {
	return &regexpParser{
		re:   regexp,
		data: make(GroupMap),
	}
}

type regexpParser struct {
	re *regexp.Regexp

	data GroupMap
}

func (rp *regexpParser) Parse(line string) (GroupMap, bool) {
	lines := rp.re.FindStringSubmatch(line)

	if lines == nil {
		return nil, false
	}

	for i, v := range rp.re.SubexpNames()[1:] {
		rp.data[v] = lines[i+1]
	}

	return rp.data, true
}
