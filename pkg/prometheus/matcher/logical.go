package matcher

import (
	"github.com/prometheus/prometheus/pkg/labels"
)

type (
	trueMatcher  struct{}
	falseMatcher struct{}
	negMatcher   struct{ m Matcher }
	andMatcher   struct{ lhs, rhs Matcher }
	orMatcher    struct{ lhs, rhs Matcher }
)

func (trueMatcher) Matches(_ labels.Labels) bool    { return true }
func (falseMatcher) Matches(_ labels.Labels) bool   { return false }
func (m negMatcher) Matches(lbs labels.Labels) bool { return !m.m.Matches(lbs) }
func (m andMatcher) Matches(lbs labels.Labels) bool { return m.lhs.Matches(lbs) && m.rhs.Matches(lbs) }
func (m orMatcher) Matches(lbs labels.Labels) bool  { return m.lhs.Matches(lbs) || m.rhs.Matches(lbs) }

func and(lhs, rhs Matcher, others ...Matcher) andMatcher {
	m := andMatcher{lhs: lhs, rhs: rhs}
	if len(others) == 0 {
		return m
	}
	return and(m, others[0], others[1:]...)
}

func or(lhs, rhs Matcher, others ...Matcher) orMatcher {
	m := orMatcher{lhs: lhs, rhs: rhs}
	if len(others) == 0 {
		return m
	}
	return or(m, others[0], others[1:]...)
}

func not(m Matcher) Matcher {
	return negMatcher{m}
}
