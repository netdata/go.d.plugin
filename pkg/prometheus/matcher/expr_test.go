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
		"empty (Not nil): both includes And excludes": {
			expr: Expr{
				Allow: []string{},
				Deny:  []string{},
			},
			expected: true,
		},
		"empty (nil): both includes And excludes": {
			expected: true,
		},
		"empty: only includes": {
			expr: Expr{
				Deny: []string{""},
			},
			expected: false,
		},
		"empty: only excludes": {
			expr: Expr{
				Allow: []string{""},
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
		"Not set: both includes And excludes": {
			expr: Expr{},
		},
		"set: both includes And excludes": {
			expr: Expr{
				Allow: []string{
					"go_memstats_*",
					"node_*",
				},
				Deny: []string{
					"go_memstats_frees_total",
					"node_cooling_*",
				},
			},
			expectedMatcher: andMatcher{
				lhs: orMatcher{
					lhs: mustGlobName("go_memstats_*"),
					rhs: mustGlobName("node_*"),
				},
				rhs: Not(orMatcher{
					lhs: mustGlobName("go_memstats_frees_total"),
					rhs: mustGlobName("node_cooling_*"),
				}),
			},
		},
		"set: only includes": {
			expr: Expr{
				Allow: []string{
					"go_memstats_*",
					"node_*",
				},
			},
			expectedMatcher: andMatcher{
				lhs: orMatcher{
					lhs: mustGlobName("go_memstats_*"),
					rhs: mustGlobName("node_*"),
				},
				rhs: Not(falseMatcher{}),
			},
		},
		"set: only excludes": {
			expr: Expr{
				Deny: []string{
					"go_memstats_frees_total",
					"node_cooling_*",
				},
			},
			expectedMatcher: andMatcher{
				lhs: trueMatcher{},
				rhs: Not(orMatcher{
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
