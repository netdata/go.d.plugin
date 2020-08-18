package matcher

import (
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrueMatcher_Matches(t *testing.T) {
	tests := map[string]struct {
		m        trueMatcher
		lbs      labels.Labels
		expected bool
	}{
		"not empty labels": {
			lbs:      labels.Labels{{Name: labels.MetricName, Value: "name"}},
			expected: true,
		},
		"empty labels": {
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.expected {
				assert.True(t, test.m.Matches(test.lbs))
			} else {
				assert.False(t, test.m.Matches(test.lbs))
			}
		})
	}
}

func TestFalseMatcher_Matches(t *testing.T) {
	tests := map[string]struct {
		m        falseMatcher
		lbs      labels.Labels
		expected bool
	}{
		"not empty labels": {
			lbs:      labels.Labels{{Name: labels.MetricName, Value: "name"}},
			expected: false,
		},
		"empty labels": {
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.expected {
				assert.True(t, test.m.Matches(test.lbs))
			} else {
				assert.False(t, test.m.Matches(test.lbs))
			}
		})
	}
}

func TestNegMatcher_Matches(t *testing.T) {
	tests := map[string]struct {
		m        negMatcher
		lbs      labels.Labels
		expected bool
	}{
		"true matcher": {
			m:        negMatcher{trueMatcher{}},
			lbs:      labels.Labels{{Name: labels.MetricName, Value: "name"}},
			expected: false,
		},
		"false matcher": {
			m:        negMatcher{falseMatcher{}},
			lbs:      labels.Labels{{Name: labels.MetricName, Value: "name"}},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.expected {
				assert.True(t, test.m.Matches(test.lbs))
			} else {
				assert.False(t, test.m.Matches(test.lbs))
			}
		})
	}
}

func TestAndMatcher_Matches(t *testing.T) {
	tests := map[string]struct {
		m        andMatcher
		lbs      labels.Labels
		expected bool
	}{
		"true, true": {
			m:        andMatcher{lhs: trueMatcher{}, rhs: trueMatcher{}},
			expected: true,
		},
		"true, false": {
			m:        andMatcher{lhs: trueMatcher{}, rhs: falseMatcher{}},
			expected: false,
		},
		"false, true": {
			m:        andMatcher{lhs: trueMatcher{}, rhs: falseMatcher{}},
			expected: false,
		},
		"false, false": {
			m:        andMatcher{lhs: falseMatcher{}, rhs: falseMatcher{}},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.m.Matches(test.lbs))
		})
	}
}

func TestOrMatcher_Matches(t *testing.T) {
	tests := map[string]struct {
		m        orMatcher
		lbs      labels.Labels
		expected bool
	}{
		"true, true": {
			m:        orMatcher{lhs: trueMatcher{}, rhs: trueMatcher{}},
			expected: true,
		},
		"true, false": {
			m:        orMatcher{lhs: trueMatcher{}, rhs: falseMatcher{}},
			expected: true,
		},
		"false, true": {
			m:        orMatcher{lhs: trueMatcher{}, rhs: falseMatcher{}},
			expected: true,
		},
		"false, false": {
			m:        orMatcher{lhs: falseMatcher{}, rhs: falseMatcher{}},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.m.Matches(test.lbs))
		})
	}
}

func Test_and(t *testing.T) {
	tests := map[string]struct {
		ms       []Matcher
		expected Matcher
	}{
		"2 matchers": {
			ms: []Matcher{trueMatcher{}, trueMatcher{}},
			expected: andMatcher{
				lhs: trueMatcher{},
				rhs: trueMatcher{},
			},
		},
		"4 matchers": {
			ms: []Matcher{trueMatcher{}, trueMatcher{}, trueMatcher{}, trueMatcher{}},
			expected: andMatcher{
				lhs: andMatcher{
					lhs: andMatcher{
						lhs: trueMatcher{},
						rhs: trueMatcher{},
					},
					rhs: trueMatcher{},
				},
				rhs: trueMatcher{}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.GreaterOrEqual(t, len(test.ms), 2)

			m := and(test.ms[0], test.ms[1], test.ms[2:]...)
			assert.Equal(t, test.expected, m)
		})
	}
}

func Test_or(t *testing.T) {
	tests := map[string]struct {
		ms       []Matcher
		expected Matcher
	}{
		"2 matchers": {
			ms: []Matcher{trueMatcher{}, trueMatcher{}},
			expected: orMatcher{
				lhs: trueMatcher{},
				rhs: trueMatcher{},
			},
		},
		"4 matchers": {
			ms: []Matcher{trueMatcher{}, trueMatcher{}, trueMatcher{}, trueMatcher{}},
			expected: orMatcher{
				lhs: orMatcher{
					lhs: orMatcher{
						lhs: trueMatcher{},
						rhs: trueMatcher{},
					},
					rhs: trueMatcher{},
				},
				rhs: trueMatcher{}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.GreaterOrEqual(t, len(test.ms), 2)

			m := or(test.ms[0], test.ms[1], test.ms[2:]...)
			assert.Equal(t, test.expected, m)
		})
	}
}
