package tcpcheck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTcpCheck_Init(t *testing.T) {
	tc := New()
	tc.Host = "127.0.0.1"
	tc.Ports = []int{123, 124}

	assert.True(t, tc.Init())

	assert.Len(t, tc.ports, 2)
	assert.Len(t, tc.workers, 2)

	time.Sleep(time.Millisecond * 200)

	for _, w := range tc.workers {
		assert.True(t, w.alive)
	}
}

func TestTcpCheck_Check(t *testing.T) {
	assert.True(t, New().Check())

}

func TestTcpCheck_Cleanup(t *testing.T) {
	tc := New()
	tc.Host = "127.0.0.1"
	tc.Ports = []int{123, 124}

	tc.Init()

	time.Sleep(time.Millisecond * 200)
	tc.Cleanup()
	time.Sleep(time.Millisecond * 200)

	for _, w := range tc.workers {
		assert.False(t, w.alive)
	}

}

func TestTcpCheck_GetCharts(t *testing.T) {
	assert.NotNil(t, New().GetCharts())

}

func TestTcpCheck_GetData(t *testing.T) {
	tc := New()
	tc.Host = "127.0.0.1"
	tc.Ports = []int{123, 124}

	tc.Init()
	time.Sleep(time.Millisecond * 200)

	assert.NotNil(t, tc.GetData())

	for _, port := range tc.ports {
		assert.Equal(t, port.state, failed)
		assert.Equal(t, port.inState, port.updateEvery)
	}

	assert.NotNil(t, tc.GetData())

	for _, port := range tc.ports {
		assert.Equal(t, port.state, failed)
		assert.Equal(t, port.inState, port.updateEvery*2)
	}
}
