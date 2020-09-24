package couchdb

import (
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}
