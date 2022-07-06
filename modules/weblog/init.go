// SPDX-License-Identifier: GPL-3.0-or-later

package weblog

import (
	"bytes"
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
		return nil, errors.New("empty 'name' or 'match'")
	}

	m, err := matcher.Parse(up.Match)
	if err != nil {
		return nil, err
	}
	return &pattern{name: up.Name, Matcher: m}, nil
}

func (w *WebLog) createURLPatterns() error {
	if len(w.URLPatterns) == 0 {
		w.Debug("skipping URL patterns creating, no patterns provided")
		return nil
	}
	w.Debug("starting URL patterns creating")
	for _, up := range w.URLPatterns {
		p, err := newPattern(up)
		if err != nil {
			return fmt.Errorf("create pattern %+v: %v", up, err)
		}
		w.Debugf("created pattern '%s', type '%T', match '%s'", p.name, p.Matcher, up.Match)
		w.urlPatterns = append(w.urlPatterns, p)
	}
	w.Debugf("created %d URL pattern(s)", len(w.URLPatterns))
	return nil
}

func (w *WebLog) createCustomFields() error {
	if len(w.CustomFields) == 0 {
		w.Debug("skipping custom fields creating, no custom fields provided")
		return nil
	}

	w.Debug("starting custom fields creating")
	w.customFields = make(map[string][]*pattern)
	for i, cf := range w.CustomFields {
		if cf.Name == "" {
			return fmt.Errorf("create custom field: name not set (field %d)", i+1)
		}
		for _, up := range cf.Patterns {
			p, err := newPattern(up)
			if err != nil {
				return fmt.Errorf("create field '%s' pattern %+v: %v", cf.Name, up, err)
			}
			w.Debugf("created field '%s', pattern '%s', type '%T', match '%s'", cf.Name, p.name, p.Matcher, up.Match)
			w.customFields[cf.Name] = append(w.customFields[cf.Name], p)
		}
	}
	w.Debugf("created %d custom field(s)", len(w.CustomFields))
	return nil
}

func (w *WebLog) createCustomTimeFields() error {
	if len(w.CustomTimeFields) == 0 {
		w.Debug("skipping custom time fields creating, no custom time fields provided")
		return nil
	}

	w.Debug("starting custom time fields creating")
	w.customTimeFields = make(map[string][]float64)
	for i, ctf := range w.CustomTimeFields {
		if ctf.Name == "" {
			return fmt.Errorf("create custom field: name not set (field %d)", i+1)
		}
		w.customTimeFields[ctf.Name] = ctf.Histogram
		w.Debugf("created time field '%s', histogram '%v'", ctf.Name, ctf.Histogram)
	}
	w.Debugf("created %d custom time field(s)", len(w.CustomTimeFields))
	return nil
}

func (w *WebLog) createLogLine() {
	w.line = newEmptyLogLine()
	for v := range w.customFields {
		w.line.custom.fields[v] = struct{}{}
	}
	for v := range w.customTimeFields {
		w.line.custom.fields[v] = struct{}{}
	}
}

func (w *WebLog) createLogReader() error {
	w.Cleanup()
	w.Debug("starting log reader creating")
	reader, err := logs.Open(w.Path, w.ExcludePath, w.Logger)
	if err != nil {
		return fmt.Errorf("creating log reader: %v", err)
	}
	w.Debugf("created log reader, current file '%s'", reader.CurrentFilename())
	w.file = reader
	return nil
}

func (w *WebLog) createParser() error {
	w.Debug("starting parser creating")
	lastLine, err := logs.ReadLastLine(w.file.CurrentFilename(), 0)
	if err != nil {
		return fmt.Errorf("read last line: %v", err)
	}
	lastLine = bytes.TrimRight(lastLine, "\n")
	w.Debugf("last line: '%s'", string(lastLine))

	w.parser, err = w.newParser(lastLine)
	if err != nil {
		return fmt.Errorf("create parser: %v", err)
	}
	w.Debugf("created parser: %s", w.parser.Info())

	err = w.parser.Parse(lastLine, w.line)
	if err != nil {
		return fmt.Errorf("parse last line: %v (%s)", err, string(lastLine))
	}

	if err = w.line.verify(); err != nil {
		return fmt.Errorf("verify last line: %v (%s)", err, string(lastLine))
	}
	return nil
}
