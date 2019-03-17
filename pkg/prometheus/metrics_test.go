package prometheus

import (
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: write better tests

const (
	testName1 = "logback_events_total"
	testName2 = "jvm_threads_peak"
)

func testMetrics() Metrics {
	return Metrics{
		{
			Value: 10,
			Labels: labels.Labels{
				{Name: "__name__", Value: testName1},
				{Name: "level", Value: "error"},
			},
		},
		{
			Value: 20,
			Labels: labels.Labels{
				{Name: "__name__", Value: testName1},
				{Name: "level", Value: "warn"},
			},
		},
		{
			Value: 5,
			Labels: labels.Labels{
				{Name: "__name__", Value: testName1},
				{Name: "level", Value: "info"},
			},
		},
		{
			Value: 15,
			Labels: labels.Labels{
				{Name: "__name__", Value: testName1},
				{Name: "level", Value: "debug"},
			},
		},
		{
			Value: 26,
			Labels: labels.Labels{
				{Name: "__name__", Value: testName2},
			},
		},
	}
}

func TestMetric_Name(t *testing.T) {
	m := testMetrics()

	assert.Equal(t, testName1, m[0].Name())
	assert.Equal(t, testName1, m[1].Name())

}

func TestMetrics_Add(t *testing.T) {
	m := testMetrics()

	require.Len(t, m, 5)
	m.Add(Metric{})
	assert.Len(t, m, 6)
}

func TestMetrics_FindByName(t *testing.T) {
	m := testMetrics()
	m.Sort()
	assert.Len(t, Metrics{}.FindByName(testName1), 0)
	assert.Len(t, m.FindByName(testName1), len(m)-1)
}

func TestMetrics_FindByNames(t *testing.T) {
	m := testMetrics()
	m.Sort()
	assert.Len(t, m.FindByNames(), 0)
	assert.Len(t, m.FindByNames(testName1), len(m)-1)
	assert.Len(t, m.FindByNames(testName1, testName2), len(m))
}

func TestMetrics_Len(t *testing.T) {
	m := testMetrics()

	assert.Equal(t, len(m), m.Len())
}

func TestMetrics_Less(t *testing.T) {
	m := testMetrics()

	assert.False(t, m.Less(0, 1))
	assert.True(t, m.Less(4, 0))
}

func TestMetrics_Match(t *testing.T) {
	m := testMetrics()

	assert.Len(
		t,
		m.Match(&labels.Matcher{
			Type:  labels.MatchEqual,
			Name:  "__name__",
			Value: testName1,
		}),
		4,
	)

}

func TestMetrics_Max(t *testing.T) {
	m := testMetrics()

	assert.Equal(t, float64(26), m.Max())

}

func TestMetrics_Reset(t *testing.T) {
	m := testMetrics()
	m.Reset()

	assert.Len(t, m, 0)

}

func TestMetrics_Sort(t *testing.T) {
	{
		m := testMetrics()
		m.Sort()
		assert.Equal(t, testName2, m[0].Name())
	}
	{
		m := Metrics{}
		assert.Equal(t, 0.0, m.Max())
	}
}

func TestMetrics_Swap(t *testing.T) {
	m := testMetrics()

	m0 := m[0]
	m1 := m[1]

	m.Swap(0, 1)

	assert.Equal(t, m0, m[1])
	assert.Equal(t, m1, m[0])
}
