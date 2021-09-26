package socket

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testServerAddress     = "127.0.0.1:9999"
	testUnixServerAddress = "/tmp/testSocketFD"
)

var tcpConfig = Config{
	Network: "tcp",
	Address: testServerAddress,
	Timeout: 5 * time.Millisecond,
	TlsConf: nil,
}

var udpConfig = Config{
	Network: "udp",
	Address: testServerAddress,
	Timeout: 200 * time.Millisecond,
	TlsConf: nil,
}

var unixConfig = Config{
	Network: "unix",
	Address: testUnixServerAddress,
	Timeout: 2000 * time.Millisecond,
	TlsConf: nil,
}

func Test_clientCommand(t *testing.T) {
	srv := &tcpServer{addr: testServerAddress, rowsNumResp: 1}
	go func() { _ = srv.Run() }()
	defer func() { _ = srv.Close() }()
	time.Sleep(time.Millisecond * 300)
	sock := NewSocket(tcpConfig)
	require.NoError(t, sock.Connect())
	err := sock.Command("ping\n", func(bytes []byte) bool {
		assert.Equal(t, "pong", string(bytes))
		return true
	})
	require.NoError(t, sock.Disconnect())
	require.NoError(t, err)
}

func Test_clientCommandStopProcessing(t *testing.T) {
	srv := &tcpServer{addr: testServerAddress, rowsNumResp: 2}
	go func() { _ = srv.Run() }()
	defer func() { _ = srv.Close() }()
	time.Sleep(time.Millisecond * 300)
	sock := NewSocket(tcpConfig)
	require.NoError(t, sock.Connect())
	err := sock.Command("ping\n", func(bytes []byte) bool {
		assert.Equal(t, "pong", string(bytes))
		return false
	})
	require.NoError(t, sock.Disconnect())
	require.NoError(t, err)
}

func Test_clientUDPCommand(t *testing.T) {
	srv := &udpServer{addr: testServerAddress, rowsNumResp: 1}
	go func() { _ = srv.Run() }()
	defer func() { _ = srv.Close() }()
	time.Sleep(time.Millisecond * 500)
	sock := NewSocket(udpConfig)
	require.NoError(t, sock.Connect())
	err := sock.Command("ping\n", func(bytes []byte) bool {
		assert.Equal(t, "pong", string(bytes))
		return false
	})
	require.NoError(t, sock.Disconnect())
	require.NoError(t, err)
}

func Test_clientUnixCommand(t *testing.T) {
	srv := &unixServer{addr: testUnixServerAddress, rowsNumResp: 1}
	go func() { _ = srv.Run() }()
	defer func() { _ = srv.Close() }()
	time.Sleep(time.Millisecond * 1000)
	sock := NewSocket(unixConfig)
	require.NoError(t, sock.Connect())
	err := sock.Command("ping\n", func(bytes []byte) bool {
		assert.Equal(t, "pong", string(bytes))
		return false
	})
	require.NoError(t, sock.Disconnect())
	require.NoError(t, err)
}

type tcpServer struct {
	addr        string
	server      net.Listener
	rowsNumResp int
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
			return errors.New("could not accept connection")
		}
		go t.handleConnection(conn)
	}
}

func (t *tcpServer) handleConnection(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	_ = conn.SetDeadline(time.Now().Add(time.Second))

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	_, err := rw.ReadString('\n')
	if err != nil {
		_, _ = rw.WriteString("failed to read input")
		_ = rw.Flush()
	} else {
		resp := strings.Repeat("pong\n", t.rowsNumResp)
		_, _ = rw.WriteString(resp)
		_ = rw.Flush()
	}
}

type udpServer struct {
	addr        string
	conn        *net.UDPConn
	rowsNumResp int
}

func (u *udpServer) Run() (err error) {
	addr, err := net.ResolveUDPAddr("udp", u.addr)
	if err != nil {
		return err
	}
	u.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	go u.handleConnections()
	return nil
}

func (u *udpServer) Close() (err error) {
	return u.conn.Close()
}

func (u *udpServer) handleConnections() {
	for {
		var buf [2048]byte
		_, addr, _ := u.conn.ReadFromUDP(buf[0:])
		resp := strings.Repeat("pong\n", u.rowsNumResp)
		_, _ = u.conn.WriteToUDP([]byte(resp), addr)
	}
}

type unixServer struct {
	addr        string
	conn        *net.UnixListener
	rowsNumResp int
}

func (u *unixServer) Run() (err error) {
	_, _ = ioutil.TempFile("/tmp", "testSocketFD")
	addr, err := net.ResolveUnixAddr("unix", u.addr)
	if err != nil {
		return err
	}
	u.conn, err = net.ListenUnix("unix", addr)
	if err != nil {
		return
	}
	go u.handleConnections()
	return nil
}

func (u *unixServer) Close() (err error) {
	_ = os.Remove(testUnixServerAddress)
	return u.conn.Close()
}

func (u *unixServer) handleConnections() {
	var conn net.Conn
	var err error
	conn, err = u.conn.AcceptUnix()
	if err != nil || conn == nil {
		panic(fmt.Errorf("could not accept connection: %v", err))
	}
	u.handleConnection(conn)
}

func (u *unixServer) handleConnection(conn net.Conn) {
	_ = conn.SetDeadline(time.Now().Add(time.Second))

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	_, err := rw.ReadString('\n')
	if err != nil {
		_, _ = rw.WriteString("failed to read input")
		_ = rw.Flush()
	} else {
		resp := strings.Repeat("pong\n", u.rowsNumResp)
		_, _ = rw.WriteString(resp)
		_ = rw.Flush()
	}
}
