package multipath

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Len(
		t,
		New("path1", "path2", "path2", "", "path3"),
		3,
	)
}

func TestMultiPath_Find(t *testing.T) {
	m := New("path1", "testdata")

	v, err := m.Find("not exist")
	assert.Zero(t, v)
	assert.Error(t, err)

	v, err = m.Find("test-empty.conf")
	assert.Equal(t, "testdata/test-empty.conf", v)
	assert.Nil(t, err)

	v, err = m.Find("test.conf")
	assert.Equal(t, "testdata/test.conf", v)
	assert.Nil(t, err)
}

func TestIsNotFound(t *testing.T) {
	assert.True(t, IsNotFound(ErrNotFound{}))
	assert.False(t, IsNotFound(errors.New("")))
}
