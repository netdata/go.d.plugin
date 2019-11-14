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

const (
	srvAddress = "127.0.0.1:38001"
)

func Test_clientSend(t *testing.T) {
	srv := &tcpServer{addr: srvAddress, respNumLines: 10}
	go srv.Run()
	defer srv.Close()
	time.Sleep(time.Second)

	c := newClient(clientConfig{
		address: srvAddress,
		timeout: time.Second,
	})

	lines, err := c.send("whatever\n")
	assert.NoError(t, err)
	assert.Len(t, lines, 10)

	lines, err = c.send("whatever\n")
	assert.NoError(t, err)
	assert.Len(t, lines, 10)
}

func Test_clientSend_ReadLineLimitExceeded(t *testing.T) {
	srv := &tcpServer{addr: srvAddress, respNumLines: maxLinesToRead + 1}
	go srv.Run()
	defer srv.Close()
	time.Sleep(time.Second)

	c := newClient(clientConfig{
		address: srvAddress,
		timeout: time.Second,
	})

	lines, err := c.send("whatever\n")
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
