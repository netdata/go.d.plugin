package portcheck

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortCheck_Init(t *testing.T) {
	tc := New()
	defer tc.Cleanup()

	tc.Host = "127.0.0.1"
	tc.Ports = []int{3001, 3002}

	assert.True(t, tc.Init())

	assert.Len(t, tc.ports, 2)
	assert.Len(t, tc.workers, 2)

	for _, w := range tc.workers {
		assert.True(t, w.alive)
	}
}

func TestPortCheck_Check(t *testing.T) {
	tc := New()
	defer tc.Cleanup()

	assert.True(t, tc.Check())
}

func TestPortCheck_Cleanup(t *testing.T) {
	tc := New()
	tc.Host = "127.0.0.1"
	tc.Ports = []int{3001, 3002}

	tc.Init()
	tc.Cleanup()

	for _, w := range tc.workers {
		assert.False(t, w.alive)
	}

}

func TestPortCheck_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestPortCheck_ServerOK(t *testing.T) {
	tc := New()
	defer tc.Cleanup()

	tc.Host = "127.0.0.1"
	tc.Ports = []int{3001}

	tc.Init()

	srv := tcpServer{addr: ":3001"}
	_ = srv.listen()

	defer func() {
		_ = srv.close()
	}()

	assert.NotNil(t, tc.GatherMetrics())

	for _, p := range tc.ports {
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
