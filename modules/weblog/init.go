package weblog

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/logs"
	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type category struct {
	name    string
	Matcher matcher.Matcher
}

func newCategory(raw rawCategory) (*category, error) {
	if raw.Name == "" || raw.Match == "" {
		return nil, fmt.Errorf("category bad syntax")
	}

	m, err := matcher.Parse(raw.Match)
	if err != nil {
		return nil, err
	}

	return &category{name: raw.Name, Matcher: m}, nil
}

func (w *WebLog) initFilter() (err error) {
	if w.Filter.Empty() {
		w.filter = matcher.TRUE()
		return
	}

	m, err := w.Filter.Parse()
	if err != nil {
		return fmt.Errorf("error on creating filter %s: %v", w.Filter, err)
	}

	w.filter = m
	return
}

func (w *WebLog) initCategories() error {
	for _, raw := range w.URLCategories {
		cat, err := newCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating url categories %s: %v", raw, err)
		}
		w.urlCategories = append(w.urlCategories, cat)
	}

	for _, raw := range w.UserCategories {
		cat, err := newCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating user categories %s: %v", raw, err)
		}
		w.userCategories = append(w.userCategories, cat)
	}

	return nil
}

func (w *WebLog) initLogReader() error {
	file, err := logs.Open(w.Path, w.ExcludePath, w.Logger)
	if err != nil {
		return fmt.Errorf("error on creating logreader : %v", err)
	}

	w.file = file
	return nil
}

func (w *WebLog) initParser() error {
	lastLine, err := logs.ReadLastLine(w.file.CurrentFilename(), 0)
	if err != nil {
		return fmt.Errorf("error on reading last line : %v", err)
	}

	w.parser, err = newParser(w.Config.Parser, w.file, lastLine)
	if err != nil {
		return fmt.Errorf("error on creating parser : %v", err)
	}

	logLine := newEmptyLogLine()
	err = w.parser.Parse(lastLine, logLine)
	if err != nil {
		return fmt.Errorf("error on parsing last line : %v (%s)", err, lastLine)
	}

	if err = logLine.Verify(); err != nil {
		return fmt.Errorf("error on verifying parsed log line : %v", err)
	}

	w.line = logLine
	return nil
}
