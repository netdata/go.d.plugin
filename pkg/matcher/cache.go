package matcher

import (
	"github.com/hashicorp/golang-lru"
)

type (
	cachedMatcher struct {
		matcher Matcher
		cache   *lru.Cache
	}
)

// WithCache adds limited cache to the Matcher.
// Limit must be a positive number. Otherwise WithCache will panic.
func WithCache(m Matcher, limit int) Matcher {
	switch m {
	case TRUE(), FALSE():
		return m
	default:
		cm := &cachedMatcher{matcher: m}
		cache, err := lru.New(limit)
		if err != nil {
			panic(err)
		}
		cm.cache = cache
		return cm
	}
}

func (m *cachedMatcher) Match(b []byte) bool {
	s := string(b)
	if m.cache == nil {
		return m.matcher.MatchString(s)
	}

	if result, ok := m.fetch(s); ok {
		return result
	}
	result := m.matcher.Match(b)
	m.put(s, result)
	return result
}

func (m *cachedMatcher) MatchString(s string) bool {
	if m.cache == nil {
		return m.matcher.MatchString(s)
	}

	if result, ok := m.fetch(s); ok {
		return result
	}
	result := m.matcher.MatchString(s)
	m.put(s, result)
	return result
}

func (m *cachedMatcher) fetch(key string) (result bool, exist bool) {
	var v interface{}
	v, exist = m.cache.Get(key)
	if !exist {
		return
	}

	result = v.(bool)
	return
}

func (m *cachedMatcher) put(key string, result bool) {
	m.cache.Add(key, result)
}
