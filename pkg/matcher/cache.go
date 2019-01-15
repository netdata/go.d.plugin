package matcher

import "sync"

type (
	cachedMatcher struct {
		matcher Matcher

		limit   int
		mtx     sync.RWMutex
		cache   map[string]bool
		inCache int
	}
)

// WithCache adds limited cache to the Matcher.
// Limit <=0 means no limit. Cache is reset after reaching the limit.
func WithCache(m Matcher, limit int) Matcher {
	switch m {
	case TRUE(), FALSE():
		return m
	default:
		return &cachedMatcher{
			limit:   limit,
			matcher: m,
			cache:   map[string]bool{},
		}
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
	m.mtx.RLock()
	result, ok = m.cache[key]
	m.mtx.RUnlock()
	return
}

func (m *cachedMatcher) put(key string, result bool) {
	m.mtx.Lock()
	if m.limit > 0 && m.inCache >= m.limit {
		m.cache = make(map[string]bool)
		m.inCache = 0
	}
	m.cache[key] = result
	m.inCache++
	m.mtx.Unlock()
}
