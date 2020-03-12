package unbound

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_clientNewNetwork(t *testing.T) {
	tests := []struct {
		address     string
		wantNetwork string
	}{
		{"127.0.0.1", "tcp"},
		{"127.0.0.1:8953", "tcp"},
		{"/usr/local/etc/unbound/run/unbound.ctl", "unix"},
	}

	for _, tt := range tests {
		t.Run(tt.address, func(t *testing.T) {
			conf := clientConfig{address: tt.address}
			cl := newClient(conf)
			assert.Equalf(t, tt.wantNetwork, cl.network, "expected '%s' client network, got '%s'", tt.wantNetwork, cl.network)
		})
	}
}

func Test_clientSend(t *testing.T) {
	const numLines = 10
	cl, srv := prepareClientServer(numLines)
	defer srv.Close()
	time.Sleep(time.Second)

	lines, err := cl.send("whatever\n")
	assert.NoError(t, err)
	assert.Len(t, lines, numLines)

	lines, err = cl.send("whatever\n")
	assert.NoError(t, err)
	assert.Len(t, lines, numLines)
}

func Test_clientSend_ReadLineLimitExceeded(t *testing.T) {
	cl, srv := prepareClientServer(maxLinesToRead + 1)
	defer srv.Close()
	time.Sleep(time.Second)

	lines, err := cl.send("whatever\n")
	assert.Error(t, err)
	assert.Len(t, lines, 0)
}

type tcpServer struct {
	addr         string
	server       net.Listener
	respNumLines int
}

func (t *tcpServer) Run() (err error) {
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		return
	}
	return t.handleConnections()
}

func (t *tcpServer) Close() (err error) {
	return t.server.Close()
}

func (t *tcpServer) handleConnections() error {
	for {
		conn, err := t.server.Accept()
		if err != nil || conn == nil {
			return errors.New("could not accept connection")
		}
		go t.handleConnection(conn)
	}
}

func (t *tcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(time.Second * 2))

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	req, err := rw.ReadString('\n')
	if err != nil {
		_, _ = rw.WriteString("error failed to read input")
		_ = rw.Flush()
		return
	}

	resp := strings.Repeat(req, t.respNumLines)
	_, _ = rw.WriteString(resp)
	_ = rw.Flush()
}

func prepareClientServer(respNumLines int) (*client, *tcpServer) {
	srv := &tcpServer{
		addr:         "127.0.0.1:38002",
		respNumLines: respNumLines,
	}
	go func() { _ = srv.Run() }()

	cl := newClient(clientConfig{
		address: srv.addr,
		timeout: time.Second,
	})
	return cl, srv
}
