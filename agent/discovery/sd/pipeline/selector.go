// SPDX-License-Identifier: GPL-3.0-or-later

package pipeline

import (
	"errors"
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
)

type Selector interface {
	Matches(model.Tags) bool
}

type (
	exactSelector string
	trueSelector  struct{}
	negSelector   struct{ Selector }
	orSelector    struct{ lhs, rhs Selector }
	andSelector   struct{ lhs, rhs Selector }
)

func (s exactSelector) Matches(tags model.Tags) bool { _, ok := tags[string(s)]; return ok }
func (s trueSelector) Matches(model.Tags) bool       { return true }
func (s negSelector) Matches(tags model.Tags) bool   { return !s.Selector.Matches(tags) }
func (s orSelector) Matches(tags model.Tags) bool    { return s.lhs.Matches(tags) || s.rhs.Matches(tags) }
func (s andSelector) Matches(tags model.Tags) bool   { return s.lhs.Matches(tags) && s.rhs.Matches(tags) }

func (s exactSelector) String() string { return "{" + string(s) + "}" }
func (s negSelector) String() string   { return "{!" + stringify(s.Selector) + "}" }
func (s trueSelector) String() string  { return "{*}" }
func (s orSelector) String() string    { return "{" + stringify(s.lhs) + "|" + stringify(s.rhs) + "}" }
func (s andSelector) String() string   { return "{" + stringify(s.lhs) + ", " + stringify(s.rhs) + "}" }
func stringify(sr Selector) string     { return strings.Trim(fmt.Sprintf("%s", sr), "{}") }

func ParseSelector(line string) (sr Selector, err error) {
	words := strings.Fields(line)
	if len(words) == 0 {
		return trueSelector{}, nil
	}

	var srs []Selector
	for _, word := range words {
		if idx := strings.IndexByte(word, '|'); idx > 0 {
			sr, err = parseOrSelectorWord(word)
		} else {
			sr, err = parseSingleSelectorWord(word)
		}
		if err != nil {
			return nil, fmt.Errorf("selector '%s' contains selector '%s' with forbidden symbol", line, word)
		}
		srs = append(srs, sr)
	}

	switch len(srs) {
	case 0:
		return trueSelector{}, nil
	case 1:
		return srs[0], nil
	default:
		return newAndSelector(srs[0], srs[1], srs[2:]...), nil
	}
}

func MustParseSelector(line string) Selector {
	sr, err := ParseSelector(line)
	if err != nil {
		panic(fmt.Sprintf("selector '%s' parse error: %v", line, err))
	}
	return sr
}

func parseOrSelectorWord(orWord string) (sr Selector, err error) {
	var srs []Selector
	for _, word := range strings.Split(orWord, "|") {
		if sr, err = parseSingleSelectorWord(word); err != nil {
			return nil, err
		}
		srs = append(srs, sr)
	}
	switch len(srs) {
	case 0:
		return trueSelector{}, nil
	case 1:
		return srs[0], nil
	default:
		return newOrSelector(srs[0], srs[1], srs[2:]...), nil
	}
}

func parseSingleSelectorWord(word string) (Selector, error) {
	if len(word) == 0 {
		return nil, errors.New("empty word")
	}
	neg := word[0] == '!'
	if neg {
		word = word[1:]
	}
	if len(word) == 0 {
		return nil, errors.New("empty word")
	}
	if word != "*" && !isSelectorWordValid(word) {
		return nil, errors.New("forbidden symbol")
	}

	var sr Selector
	switch word {
	case "*":
		sr = trueSelector{}
	default:
		sr = exactSelector(word)
	}
	if neg {
		return negSelector{sr}, nil
	}
	return sr, nil
}

func newAndSelector(lhs, rhs Selector, others ...Selector) Selector {
	m := andSelector{lhs: lhs, rhs: rhs}
	switch len(others) {
	case 0:
		return m
	default:
		return newAndSelector(m, others[0], others[1:]...)
	}
}

func newOrSelector(lhs, rhs Selector, others ...Selector) Selector {
	m := orSelector{lhs: lhs, rhs: rhs}
	switch len(others) {
	case 0:
		return m
	default:
		return newOrSelector(m, others[0], others[1:]...)
	}
}

func isSelectorWordValid(word string) bool {
	// valid:
	// *
	// ^[a-zA-Z][a-zA-Z0-9=_.]*$
	if len(word) == 0 {
		return false
	}
	if word == "*" {
		return true
	}
	for i, b := range word {
		switch {
		default:
			return false
		case b >= 'a' && b <= 'z':
		case b >= 'A' && b <= 'Z':
		case b >= '0' && b <= '9' && i > 0:
		case (b == '=' || b == '_' || b == '.') && i > 0:
		}
	}
	return true
}
