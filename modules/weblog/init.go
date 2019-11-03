package weblog

import (
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
		return nil, fmt.Errorf("pattern bad syntax: %v", up)
	}

	m, err := matcher.Parse(up.Match)
	if err != nil {
		return nil, err
	}
	return &pattern{name: up.Name, Matcher: m}, nil
}

func (w *WebLog) initFilter() (err error) {
	if w.Filter.Empty() {
		w.filter = matcher.TRUE()
		return
	}

	w.filter, err = w.Filter.Parse()
	if err != nil {
		return fmt.Errorf("error on creating filter %s: %v", w.Filter, err)
	}
	return err
}

func (w *WebLog) initPatterns() error {
	for _, up := range w.URLPatterns {
		p, err := newPattern(up)
		if err != nil {
			return fmt.Errorf("error on creating url pattern %s: %v", up, err)
		}
		w.patURL = append(w.patURL, p)
	}

	for _, up := range w.CustomPatterns {
		p, err := newPattern(up)
		if err != nil {
			return fmt.Errorf("error on creating user pattern %s: %v", up, err)
		}
		w.patCustom = append(w.patCustom, p)
	}
	return nil
}

func (w *WebLog) initLogReader() error {
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

	line := newEmptyLogLine()
	err = w.parser.Parse(lastLine, line)
	if err != nil {
		return fmt.Errorf("error on parsing last line: %v (%s)", err, lastLine)
	}

	if err = line.verify(); err != nil {
		return fmt.Errorf("error on verifying parsed log line: %v", err)
	}
	w.line = line
	return nil
}
