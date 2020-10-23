package wireguard

import (
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"

	"golang.zx2c4.com/wireguard/wgctrl"
)

func isPermittedOperation() bool {
	connection, err := wgctrl.New()
	if err != nil {
		return false
	}
	_, err = connection.Devices()
	if err != nil {
		return false
	}
	return true
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestWireguard_Init(t *testing.T) {
	mod := New()
	if isPermittedOperation() {
		assert.True(t, mod.Init())
	} else {
		assert.False(t, mod.Init())
	}
}

func TestWireguard_Check(t *testing.T) {
	mod := New()

	assert.True(t, mod.Check())
}

func TestWireguard_Charts(t *testing.T) {
	mod := New()
	if isPermittedOperation() {
		assert.NotNil(t, mod.Charts())
	} else {
		assert.False(t, mod.Init())
	}
}

func TestWireguard_Cleanup(t *testing.T) {
	mod := New()

	if isPermittedOperation() {
		mod.Cleanup()
		assert.Nil(t, mod.connection)
	} else {
		assert.False(t, mod.Init())
	}
}

func TestWireguard_Collect(t *testing.T) {
	mod := New()

	if isPermittedOperation() {
		assert.NotNil(t, mod.Collect())
	} else {
		assert.False(t, mod.Init())
	}
}
