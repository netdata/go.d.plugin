package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCache(t *testing.T) {
	regMatcher, _ := NewRegExpMatcher("[0-9]+")
	cached := WithCache(regMatcher)

	assert.True(t, cached.MatchString("1"))
	assert.True(t, cached.MatchString("1"))
	assert.True(t, cached.Match([]byte("2")))
	assert.True(t, cached.Match([]byte("2")))
}

func TestWithCache_specialCase(t *testing.T) {
	assert.Equal(t, TRUE(), WithCache(TRUE()))
	assert.Equal(t, FALSE(), WithCache(FALSE()))
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
		m := WithCache(globMatcher("abc*def*ghi"))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.MatchString("abc123def456ghi")
		}
	})
}
