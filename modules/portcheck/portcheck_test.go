package portcheck

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*PortCheck)(nil), New())
}

func TestPortCheck_Init(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Host = "127.0.0.1"
	mod.Ports = []int{3001, 3002}

	assert.True(t, mod.Init())

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
	mod.Ports = []int{3001, 3002}

	mod.Init()
	mod.Cleanup()

	for _, w := range mod.workers {
		assert.False(t, w.alive)
	}

}

func TestPortCheck_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestPortCheck_GatherMetrics(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Host = "127.0.0.1"
	mod.Ports = []int{3001, 3002}

	mod.UpdateEvery = 5
	mod.Init()

	srv := tcpServer{addr: ":3001"}
	_ = srv.listen()

	defer srv.close()

	expected := map[string]int64{
		"success_3001": 1,
		"failed_3001":  0,
		"timeout_3001": 0,
		"instate_3001": 5,
		"success_3002": 0,
		"failed_3002":  1,
		"timeout_3002": 0,
		"instate_3002": 5,
	}

	rv := mod.GatherMetrics()

	require.NotNil(t, rv)

	delete(rv, "latency_3001")
	delete(rv, "latency_3002")

	assert.Equal(t, expected, rv)
}

func TestPortCheck_ServerOK(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Host = "127.0.0.1"
	mod.Ports = []int{3001}

	mod.Init()

	srv := tcpServer{addr: ":3001"}
	_ = srv.listen()

	defer srv.close()

	assert.NotNil(t, mod.GatherMetrics())

	for _, p := range mod.ports {
		assert.True(t, p.state == success)
	}
}

func TestPortCheck_ServerBAD(t *testing.T) {
	tc := New()

	defer tc.Cleanup()

	tc.Host = "127.0.0.1"
	tc.Ports = []int{3001}

	tc.Init()

	assert.NotNil(t, tc.GatherMetrics())

	for _, p := range tc.ports {
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
