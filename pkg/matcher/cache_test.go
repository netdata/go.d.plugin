package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCache(t *testing.T) {
	regMatcher, _ := NewRegExpMatcher("[0-9]+")
	cached := WithCache(regMatcher, -1)

	assert.True(t, cached.MatchString("1"))
	assert.True(t, cached.MatchString("1"))
	assert.True(t, cached.Match([]byte("2")))
	assert.True(t, cached.Match([]byte("2")))

	regMatcher, _ = NewRegExpMatcher("[0-9]+")

	cached = WithCache(regMatcher, -1)
	cm := cached.(*cachedMatcher)
	assert.IsType(t, (simpleCache)(nil), cm.cache)

	cached = WithCache(regMatcher, 4)
	cm = cached.(*cachedMatcher)
	assert.IsType(t, (*lruCache)(nil), cm.cache)
}

func TestWithCache_specialCase(t *testing.T) {
	assert.Equal(t, TRUE(), WithCache(TRUE(), -1))
	assert.Equal(t, FALSE(), WithCache(FALSE(), -1))
}

func BenchmarkCachedMatcher_Match(b *testing.B) {
	b.Run("raw", func(b *testing.B) {
		m := globMatcher("abc*def*ghi")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.MatchString("abc123def456ghi")
		}
	})
	b.Run("cached", func(b *testing.B) {
		m := WithCache(globMatcher("abc*def*ghi"), -1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.MatchString("abc123def456ghi")
		}
	})
}
