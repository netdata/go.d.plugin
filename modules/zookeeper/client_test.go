package zookeeper

import (
	"bufio"
	"errors"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
	"time"
)

const (
	testServerAddress = "127.0.0.1:38001"
)

func Test_clientFetch(t *testing.T) {
	srv := &tcpServer{addr: testServerAddress}
	go srv.Run()
	defer srv.Close()
	time.Sleep(time.Second)

	c := newClient(clientConfig{
		network: "tcp",
		address: testServerAddress,
		timeout: time.Second,
	})

	rows, err := c.fetch("whatever\n")
	assert.NoError(t, err)
	assert.Len(t, rows, 10)

	rows, err = c.fetch("whatever\n")
	assert.NoError(t, err)
	assert.Len(t, rows, 10)
}

type tcpServer struct {
	addr   string
	server net.Listener
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

func (t *tcpServer) handleConnections() (err error) {
	for {
		conn, err := t.server.Accept()
		if err != nil || conn == nil {
			err = errors.New("could not accept connection")
			break
		}

		go t.handleConnection(conn)
	}
	return
}

func (t *tcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(time.Second * 2))

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		req, err := rw.ReadString('\n')
		if err != nil {
			_, _ = rw.WriteString("failed to read input")
			_ = rw.Flush()
			return
		}

		resp := strings.Repeat(req, 10)
		_, _ = rw.WriteString(resp)
		_ = rw.Flush()
		return
	}
}
