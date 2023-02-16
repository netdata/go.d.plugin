// SPDX-License-Identifier: GPL-3.0-or-later

package vnode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	assert.NotNil(t, NewRegistry("testdata"))
	assert.NotNil(t, NewRegistry("not_exist"))
}

func TestRegistry_Lookup(t *testing.T) {
	req := NewRegistry("testdata")

	_, ok := req.Lookup("first")
	assert.True(t, ok)

	_, ok = req.Lookup("second")
	assert.True(t, ok)

	_, ok = req.Lookup("third")
	assert.False(t, ok)
}
