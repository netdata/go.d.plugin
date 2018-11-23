package portcheck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPortCheck_Init(t *testing.T) {
	pc := New()
	pc.Host = "127.0.0.1"
	pc.Ports = []int{123, 124}

	assert.True(t, pc.Init())

	assert.Len(t, pc.ports, 2)
	assert.Len(t, pc.workers, 2)

	time.Sleep(time.Millisecond * 200)

	for _, w := range pc.workers {
		assert.True(t, w.alive)
	}
}

func TestPortCheck_Check(t *testing.T) {
	assert.True(t, New().Check())

}

func TestPortCheck_Cleanup(t *testing.T) {
	pc := New()
	pc.Host = "127.0.0.1"
	pc.Ports = []int{123, 124}

	pc.Init()

	time.Sleep(time.Millisecond * 200)
	pc.Cleanup()
	time.Sleep(time.Millisecond * 200)

	for _, w := range pc.workers {
		assert.False(t, w.alive)
	}

}

func TestPortCheck_GetCharts(t *testing.T) {
	assert.NotNil(t, New().GetCharts())

}

func TestPortCheck_GetData(t *testing.T) {
	pc := New()
	pc.Host = "127.0.0.1"
	pc.Ports = []int{123, 124}

	pc.Init()
	time.Sleep(time.Millisecond * 200)

	assert.NotNil(t, pc.GetData())

	for _, port := range pc.ports {
		assert.Equal(t, port.state, failed)
		assert.Equal(t, port.inState, port.updateEvery)
	}

	assert.NotNil(t, pc.GetData())

	for _, port := range pc.ports {
		assert.Equal(t, port.state, failed)
		assert.Equal(t, port.inState, port.updateEvery*2)
	}
}
