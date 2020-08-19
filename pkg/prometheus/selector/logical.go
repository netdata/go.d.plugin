package selector

import (
	"github.com/prometheus/prometheus/pkg/labels"
)

type (
	trueSelector  struct{}
	falseSelector struct{}
	negSelector   struct{ s Selector }
	andSelector   struct{ lhs, rhs Selector }
	orSelector    struct{ lhs, rhs Selector }
)

func (trueSelector) Matches(_ labels.Labels) bool    { return true }
func (falseSelector) Matches(_ labels.Labels) bool   { return false }
func (s negSelector) Matches(lbs labels.Labels) bool { return !s.s.Matches(lbs) }
func (s andSelector) Matches(lbs labels.Labels) bool { return s.lhs.Matches(lbs) && s.rhs.Matches(lbs) }
func (s orSelector) Matches(lbs labels.Labels) bool  { return s.lhs.Matches(lbs) || s.rhs.Matches(lbs) }

// And returns a matcher which returns true only if all of it's sub-matcher return true
func And(lhs, rhs Selector, others ...Selector) Selector {
	m := andSelector{lhs: lhs, rhs: rhs}
	if len(others) == 0 {
		return m
	}
	return And(m, others[0], others[1:]...)
}

// Or returns a matcher which returns true if any of it's sub-matcher return true
func Or(lhs, rhs Selector, others ...Selector) Selector {
	m := orSelector{lhs: lhs, rhs: rhs}
	if len(others) == 0 {
		return m
	}
	return Or(m, others[0], others[1:]...)
}

// Not returns a matcher which opposites the sub-matcher's result
func Not(m Selector) Selector {
	return negSelector{m}
}
