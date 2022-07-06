// SPDX-License-Identifier: GPL-3.0-or-later

package chrony

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// We can't fake a reply based on the current request,
// otherwise we don't know if chrony's reply can really be processed,
// so all of these test need chrony running, and listen at default port.

func TestNew(t *testing.T) {
	assert.IsType(t, (*Chrony)(nil), New())
}

func TestChrony_Init(t *testing.T) {
	assert.True(t, New().Init())
}

func TestChrony_Check(t *testing.T) {
	mod := New()
	mod.Init()
	assert.True(t, mod.Check())
}

func TestChrony_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestChrony_Cleanup(t *testing.T) {
	New().Cleanup()
}

//func TestChrony_Collect(t *testing.T) {
//	mod := New()
//	mod.Init()
//
//	ans := mod.Collect()
//
//	// should have something in result
//	assert.NotNil(t, mod.Collect())
//	// chrony should be running
//	if ans["running"] == 1 {
//		// in most cases, the leap second status should be 0
//		assert.EqualValues(t, 0, ans["leap_status"])
//
//		// should collect source server
//		assert.True(t, mod.Charts().Has("source"))
//		// if chrony syncs upstream normally, the source should not be 0.0.0.0
//		assert.False(t, mod.Charts().Get("source").HasDim(net.IPv4zero.String()))
//		// if chrony syncs upstream normally, should at least one online source
//		assert.NotEqualValues(t, 0, ans["online_sources"])
//	}
//
//}
