package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// NG case
	_, err := New("")
	assert.Error(t, err)

	_, err = New("invalid=method")
	assert.Error(t, err)

	_, err = New("string=")
	assert.Error(t, err)

	m, err := New("string=expr")
	assert.NoError(t, err)
	assert.IsType(t, (*StringContainsMatcher)(nil), m)

	// OK case
	m, err = New("string=^expr")
	assert.NoError(t, err)
	assert.IsType(t, (*StringPrefixMatcher)(nil), m)

	m, err = New("string=expr$")
	assert.NoError(t, err)
	assert.IsType(t, (*StringSuffixMatcher)(nil), m)

	m, err = New("regexp=[0-9]+")
	assert.NoError(t, err)
	assert.IsType(t, (*RegexpMatcher)(nil), m)
}
