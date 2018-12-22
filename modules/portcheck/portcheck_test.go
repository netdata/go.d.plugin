package portcheck

import (
	"net"
	"testing"

	"github.com/netdata/go.d.plugin/modules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*modules.Module)(nil), New())
}

func TestPortCheck_Init(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Host = "127.0.0.1"
	mod.Ports = []int{38001, 38002}

	require.True(t, mod.Init())

	assert.Len(t, mod.ports, 2)
	assert.Len(t, mod.workers, 2)

	for _, w := range mod.workers {
		assert.True(t, w.alive)
	}
}

func TestPortCheck_Check(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	assert.True(t, mod.Check())
}

func TestPortCheck_Cleanup(t *testing.T) {
	mod := New()
	mod.Host = "127.0.0.1"
	mod.Ports = []int{38001, 38002}

	assert.True(t, mod.Init())
	mod.Cleanup()

	for _, w := range mod.workers {
		assert.False(t, w.alive)
	}

}

func TestPortCheck_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
	assert.NoError(t, modules.CheckCharts(*New().Charts()...))
}

func TestPortCheck_Collect(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Host = "127.0.0.1"
	mod.Ports = []int{38001, 38002}

	mod.UpdateEvery = 5
	require.True(t, mod.Init())

	ts := tcpServer{addr: ":38001"}
	_ = ts.listen()

	defer ts.close()

	expected := map[string]int64{
		"success_38001": 1,
		"failed_38001":  0,
		"timeout_38001": 0,
		"instate_38001": 5,
		"success_38002": 0,
		"failed_38002":  1,
		"timeout_38002": 0,
		"instate_38002": 5,
	}

	rv := mod.Collect()

	require.NotNil(t, rv)

	delete(rv, "latency_38001")
	delete(rv, "latency_38002")

	assert.Equal(t, expected, rv)
}

func TestPortCheck_ServerOK(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Host = "127.0.0.1"
	mod.Ports = []int{38001}

	require.True(t, mod.Init())

	ts := tcpServer{addr: ":38001"}
	_ = ts.listen()

	defer ts.close()

	assert.NotNil(t, mod.Collect())

	for _, p := range mod.ports {
		assert.True(t, p.state == success)
	}
}

func TestPortCheck_ServerNG(t *testing.T) {
	mod := New()

	defer mod.Cleanup()

	mod.Host = "127.0.0.1"
	mod.Ports = []int{38001}

	require.True(t, mod.Init())

	assert.NotNil(t, mod.Collect())

	for _, p := range mod.ports {
		assert.True(t, p.state == failed)
	}

}

type tcpServer struct {
	addr   string
	server net.Listener
}

func (t *tcpServer) listen() (err error) {
	t.server, err = net.Listen("tcp", t.addr)
	return err
}

func (t *tcpServer) close() error {
	return t.server.Close()
}
