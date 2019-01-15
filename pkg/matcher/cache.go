package matcher

import (
	"sync"

	"github.com/hashicorp/golang-lru"
)

type cache interface {
	get(key string) (result bool, exist bool)
	add(key string, value bool)
}

type simpleCache map[string]bool

func (c simpleCache) get(key string) (bool, bool) {
	result, ok := c[key]
	return result, ok
}

func (c simpleCache) add(key string, value bool) { c[key] = value }

func newLRUCache(limit int) cache {
	c, err := lru.New(limit)
	if err != nil {
		panic(err)
	}
	return &lruCache{c}
}

type lruCache struct {
	*lru.Cache
}

func (c *lruCache) get(key string) (bool, bool) {
	v, ok := c.Get(key)
	if !ok {
		return false, false
	}
	result := v.(bool)
	return result, ok
}

func (c *lruCache) add(key string, value bool) { c.Add(key, value) }

type (
	cachedMatcher struct {
		matcher Matcher

		mux   sync.RWMutex
		cache cache
	}
)

// WithCache adds limited cache to the matcher.
// Limit < 0 means no limit. If limit == 0 WithCache doesn't add cache to the matcher.
func WithCache(m Matcher, limit int) Matcher {
	switch m {
	case TRUE(), FALSE():
		return m
	default:
		if limit == 0 {
			return m
		}
		cm := &cachedMatcher{matcher: m}
		if limit < 0 {
			cm.cache = make(simpleCache)
		} else {
			cm.cache = newLRUCache(limit)
		}
		return cm
	}
}

func (m *cachedMatcher) Match(b []byte) bool {
	s := string(b)
	if result, ok := m.fetch(s); ok {
		return result
	}
	result := m.matcher.Match(b)
	m.put(s, result)
	return result
}

func (m *cachedMatcher) MatchString(s string) bool {
	if result, ok := m.fetch(s); ok {
		return result
	}
	result := m.matcher.MatchString(s)
	m.put(s, result)
	return result
}

func (m *cachedMatcher) fetch(key string) (result bool, ok bool) {
	m.mux.RLock()
	result, ok = m.cache.get(key)
	m.mux.RUnlock()
	return
}

func (m *cachedMatcher) put(key string, result bool) {
	m.mux.Lock()
	m.cache.add(key, result)
	m.mux.Unlock()
}
