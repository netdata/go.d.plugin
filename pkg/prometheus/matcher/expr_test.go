package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpr_Empty(t *testing.T) {
	tests := map[string]struct {
		expr     Expr
		expected bool
	}{
		"empty (not nil): both includes and excludes": {
			expr: Expr{
				Includes: []string{},
				Excludes: []string{},
			},
			expected: true,
		},
		"empty (nil): both includes and excludes": {
			expected: true,
		},
		"empty: only includes": {
			expr: Expr{
				Excludes: []string{""},
			},
			expected: false,
		},
		"empty: only excludes": {
			expr: Expr{
				Includes: []string{""},
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.expected {
				assert.True(t, test.expr.Empty())
			} else {
				assert.False(t, test.expr.Empty())
			}
		})
	}
}

func TestExpr_Parse(t *testing.T) {
	tests := map[string]struct {
		expr            Expr
		expectedMatcher Matcher
		expectedErr     bool
	}{
		"not set: both includes and excludes": {
			expr: Expr{},
		},
		"set: both includes and excludes": {
			expr: Expr{
				Includes: []string{
					"go_memstats_*",
					"node_*",
				},
				Excludes: []string{
					"go_memstats_frees_total",
					"node_cooling_*",
				},
			},
			expectedMatcher: andMatcher{
				lhs: orMatcher{
					lhs: mustGlobName("go_memstats_*"),
					rhs: mustGlobName("node_*"),
				},
				rhs: not(orMatcher{
					lhs: mustGlobName("go_memstats_frees_total"),
					rhs: mustGlobName("node_cooling_*"),
				}),
			},
		},
		"set: only includes": {
			expr: Expr{
				Includes: []string{
					"go_memstats_*",
					"node_*",
				},
			},
			expectedMatcher: andMatcher{
				lhs: orMatcher{
					lhs: mustGlobName("go_memstats_*"),
					rhs: mustGlobName("node_*"),
				},
				rhs: not(falseMatcher{}),
			},
		},
		"set: only excludes": {
			expr: Expr{
				Excludes: []string{
					"go_memstats_frees_total",
					"node_cooling_*",
				},
			},
			expectedMatcher: andMatcher{
				lhs: trueMatcher{},
				rhs: not(orMatcher{
					lhs: mustGlobName("go_memstats_frees_total"),
					rhs: mustGlobName("node_cooling_*"),
				}),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m, err := test.expr.Parse()

			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expectedMatcher, m)
			}
		})
	}
}
