// SPDX-License-Identifier: GPL-3.0-or-later

package pipeline

import (
	"regexp"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/gobwas/glob"
)

var funcMap = func() template.FuncMap {
	custom := map[string]interface{}{
		"glob": globAny,
		"re":   regexpAny,
	}

	fm := sprig.HermeticTxtFuncMap()
	for name, fn := range custom {
		fm[name] = fn
	}

	return fm
}()

func globAny(value, pattern string, rest ...string) bool {
	switch len(rest) {
	case 0:
		return globOnce(value, pattern)
	default:
		return globOnce(value, pattern) || globAny(value, rest[0], rest[1:]...)
	}
}

func regexpAny(value, pattern string, rest ...string) bool {
	switch len(rest) {
	case 0:
		return regexpOnce(value, pattern)
	default:
		return regexpOnce(value, pattern) || regexpAny(value, rest[0], rest[1:]...)
	}
}

func globOnce(value, pattern string) bool {
	g, _ := globStore(pattern)
	return g != nil && g.Match(value)
}

func regexpOnce(value, pattern string) bool {
	r, _ := regexpStore(pattern)
	return r != nil && r.MatchString(value)
}

// TODO: cleanup?
var globStore = func() func(pattern string) (glob.Glob, error) {
	var l sync.RWMutex
	store := make(map[string]struct {
		g   glob.Glob
		err error
	})

	return func(pattern string) (glob.Glob, error) {
		if pattern == "" {
			return nil, nil
		}
		l.Lock()
		defer l.Unlock()
		entry, ok := store[pattern]
		if !ok {
			entry.g, entry.err = glob.Compile(pattern, '/')
			store[pattern] = entry
		}
		return entry.g, entry.err
	}
}()

// TODO: cleanup?
var regexpStore = func() func(pattern string) (*regexp.Regexp, error) {
	var l sync.RWMutex
	store := make(map[string]struct {
		r   *regexp.Regexp
		err error
	})

	return func(pattern string) (*regexp.Regexp, error) {
		if pattern == "" {
			return nil, nil
		}
		l.Lock()
		defer l.Unlock()
		entry, ok := store[pattern]
		if !ok {
			entry.r, entry.err = regexp.Compile(pattern)
			store[pattern] = entry
		}
		return entry.r, entry.err
	}
}()
