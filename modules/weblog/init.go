package weblog

import (
	"errors"
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/logs"
	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type pattern struct {
	name string
	matcher.Matcher
}

func newPattern(up userPattern) (*pattern, error) {
	if up.Name == "" || up.Match == "" {
		return nil, fmt.Errorf("pattern bad syntax: %+v", up)
	}

	m, err := matcher.Parse(up.Match)
	if err != nil {
		return nil, err
	}
	return &pattern{name: up.Name, Matcher: m}, nil
}

func (w *WebLog) initURLPatterns() error {
	for _, up := range w.URLPatterns {
		p, err := newPattern(up)
		if err != nil {
			return fmt.Errorf("error on creating url pattern %+v: %v", up, err)
		}
		w.urlPatterns = append(w.urlPatterns, p)
	}
	return nil
}

func (w *WebLog) initCustomFields() error {
	if len(w.CustomFields) == 0 {
		return nil
	}

	w.customFields = make(map[string][]*pattern)
	for _, cf := range w.CustomFields {
		if cf.Name == "" {
			return errors.New("error on creating custom field: name not set")
		}
		for _, up := range cf.Patterns {
			p, err := newPattern(up)
			if err != nil {
				return fmt.Errorf("error on creating custom field '%s' pattern %+v: %v", cf.Name, up, err)
			}
			w.customFields[cf.Name] = append(w.customFields[cf.Name], p)
		}
	}
	return nil
}

func (w *WebLog) initLogLine() {
	w.line = newEmptyLogLine()
	for v := range w.customFields {
		w.line.custom.fields[v] = struct{}{}
	}
}

func (w *WebLog) initLogReader() error {
	w.Cleanup()
	reader, err := logs.Open(w.Path, w.ExcludePath, w.Logger)
	if err != nil {
		return fmt.Errorf("error on creating log reader: %v", err)
	}

	w.file = reader
	return nil
}

func (w *WebLog) initParser() error {
	lastLine, err := logs.ReadLastLine(w.file.CurrentFilename(), 0)
	if err != nil {
		return fmt.Errorf("error on reading last line: %v", err)
	}

	w.parser, err = w.newParser(lastLine)
	if err != nil {
		return fmt.Errorf("error on creating parser: %v", err)
	}

	err = w.parser.Parse(lastLine, w.line)
	if err != nil {
		return fmt.Errorf("error on parsing last line: %v (%s)", err, string(lastLine))
	}

	if err = w.line.verify(); err != nil {
		return fmt.Errorf("error on verifying parsed log line: %v", err)
	}
	return nil
}
