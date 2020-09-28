package couchdb

import (
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}
