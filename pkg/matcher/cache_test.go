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
}

func TestWithLimitedCache(t *testing.T) {
	regMatcher, _ := NewRegExpMatcher("[0-9]+")

	// Unlimited
	cached := WithCache(regMatcher, -1)

	for _, s := range []string{"1", "2", "3", "4", "a", "b", "c"} {
		cached.MatchString(s)
	}

	cm := cached.(*cachedMatcher)
	assert.Len(t, cm.cache, 7)
	assert.Equal(
		t,
		map[string]bool{"1": true, "2": true, "3": true, "4": true, "a": false, "b": false, "c": false},
		cm.cache,
	)

	// Limited
	cached = WithCache(regMatcher, 4)

	for _, s := range []string{"1", "2", "3", "4", "a", "b", "c"} {
		cached.MatchString(s)
	}

	cm = cached.(*cachedMatcher)
	assert.Len(t, cm.cache, 3)
	assert.Equal(
		t,
		map[string]bool{"a": false, "b": false, "c": false},
		cm.cache,
	)
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
