package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter_WriteTo(t *testing.T) {
	g := Counter{}
	g.Inc()
	g.Inc()
	g.Inc()
	g.Add(0.14)
	m := map[string]int64{}
	g.WriteTo(m, "pi", 100, 1)
	assert.Len(t, m, 1)
	assert.EqualValues(t, 314, m["pi"])
}

func TestCounter_Inc(t *testing.T) {
	c := Counter{}
	c.Inc()
	assert.Equal(t, 1.0, c.Value())
	c.Inc()
	assert.Equal(t, 2.0, c.Value())
}

func TestCounter_Add(t *testing.T) {
	c := Counter{}
	c.Add(3.14)
	assert.InDelta(t, 3.14, c.Value(), 0.0001)
	c.Add(2)
	assert.InDelta(t, 5.14, c.Value(), 0.0001)
	assert.Panics(t, func() {
		c.Add(-1)
	})
}

func BenchmarkCounter_Add(b *testing.B) {
	benchmarks := []struct {
		name  string
		value float64
	}{
		{"int", 1},
		{"float", 3.14},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var c Counter
			for i := 0; i < b.N; i++ {
				c.Add(bm.value)
			}
		})
	}
}

func BenchmarkCounter_Inc(b *testing.B) {
	var c Counter
	for i := 0; i < b.N; i++ {
		c.Inc()
	}
}

func BenchmarkCounter_Value(b *testing.B) {
	var c Counter
	c.Inc()
	c.Add(3.14)
	for i := 0; i < b.N; i++ {
		c.Value()
	}
}

func BenchmarkCounter_WriteTo(b *testing.B) {
	var c Counter
	c.Inc()
	c.Add(3.14)
	m := map[string]int64{}
	for i := 0; i < b.N; i++ {
		c.WriteTo(m, "pi", 100, 1)
	}
}
