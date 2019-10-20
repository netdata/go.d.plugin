package logs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
)

type (
	RegExpConfig struct {
		Pattern string `yaml:"pattern"`
	}

	RegExpParser struct {
		r       *bufio.Reader
		pattern *regexp.Regexp
	}
)

func NewRegExpParser(config RegExpConfig, in io.Reader) (*RegExpParser, error) {
	if config.Pattern == "" {
		return nil, errors.New("empty regexp pattern")
	}

	pattern, err := regexp.Compile(config.Pattern)
	if err != nil {
		return nil, fmt.Errorf("error on compiling regexp pattern : %w", err)
	}

	p := &RegExpParser{
		r:       bufio.NewReader(in),
		pattern: pattern,
	}
	return p, nil
}

func (p *RegExpParser) ReadLine(logLine LogLine) error {
	s, err := p.r.ReadSlice('\n')
	if err != nil && len(s) == 0 {
		return err
	}
	return p.Parse(s, logLine)
}

func (p *RegExpParser) Parse(line []byte, logLine LogLine) error {
	match := p.pattern.FindSubmatch(line)
	if match == nil {
		return &ParseError{msg: "regexp unmatched line"}
	}

	for i, name := range p.pattern.SubexpNames() {
		if name == "" || match[i] == nil {
			continue
		}
		err := logLine.Assign(name, string(match[i]))
		if err != nil {
			return &ParseError{
				msg: fmt.Sprintf("regexp error on assigning : %v", err),
				err: err,
			}
		}
	}
	return nil
}

func (p RegExpParser) Info() string {
	return p.pattern.String()
}
