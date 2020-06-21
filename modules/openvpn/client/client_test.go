package client

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testLoadStatsData, _ = ioutil.ReadFile("testdata/load-stats.txt")
	testVersionData, _   = ioutil.ReadFile("testdata/version.txt")
	testStatus3Data, _   = ioutil.ReadFile("testdata/status3.txt")

	testDefaultTimeout = time.Second
)

func testDial(conn net.Conn) dialFunc {
	return func(_, _ string, _ time.Duration) (net.Conn, error) { return conn, nil }
}

func TestNew(t *testing.T) { assert.IsType(t, (*Client)(nil), New(Config{})) }

func TestClient_Connect(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	}()

	client := &Client{dial: testDial(clientConn)}
	assert.NoError(t, client.Connect())
	assert.True(t, client.IsConnected())
}

func TestClient_Disconnect(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	}()

	client := &Client{dial: testDial(clientConn)}
	assert.False(t, client.IsConnected())
	assert.NoError(t, client.Connect())
	assert.True(t, client.IsConnected())
	assert.NoError(t, client.Disconnect())
	assert.False(t, client.IsConnected())
}

func TestClient_IsConnected(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	}()

	client := &Client{dial: testDial(clientConn)}
	assert.False(t, client.IsConnected())
	assert.NoError(t, client.Connect())
	assert.True(t, client.IsConnected())
}

func TestClient_GetVersion(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	}()
	client := newTestTCPClient(clientConn)
	srv := newTestTCPServer(serverConn)
	go srv.serve()

	ver, err := client.Version()
	assert.NoError(t, err)
	expected := &Version{Major: 2, Minor: 3, Patch: 4, Management: 1}
	assert.Equal(t, expected, ver)
}

func TestClient_GetLoadStats(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	}()
	client := newTestTCPClient(clientConn)
	srv := newTestTCPServer(serverConn)
	go srv.serve()

	stats, err := client.LoadStats()
	assert.NoError(t, err)
	expected := &LoadStats{NumOfClients: 1, BytesIn: 7811, BytesOut: 7667}
	assert.Equal(t, expected, stats)
}

func TestClient_GetUsers(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	}()
	client := newTestTCPClient(clientConn)
	srv := newTestTCPServer(serverConn)
	go srv.serve()

	users, err := client.Users()
	assert.NoError(t, err)
	expected := Users{{
		CommonName:     "pepehome",
		RealAddress:    "1.2.3.4:44347",
		VirtualAddress: "10.9.0.5",
		BytesReceived:  6043,
		BytesSent:      5661,
		ConnectedSince: 1555439465,
		Username:       "pepe",
	}}
	assert.Equal(t, expected, users)
}

func newTestTCPClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
		Config: Config{
			ConnectTimeout: testDefaultTimeout,
			ReadTimeout:    testDefaultTimeout,
			WriteTimeout:   testDefaultTimeout,
		},
	}
}

func newTestTCPServer(conn net.Conn) *testTCPServer { return &testTCPServer{conn: conn} }

type testTCPServer struct{ conn net.Conn }

func (t *testTCPServer) serve() {
	for t.serveOnce() == nil {
	}
}

func (t *testTCPServer) serveOnce() error {
	if err := t.conn.SetReadDeadline(time.Now().Add(testDefaultTimeout)); err != nil {
		return err
	}

	command, err := bufio.NewReader(t.conn).ReadString('\n')
	if err != nil {
		return err
	}

	if err = t.conn.SetWriteDeadline(time.Now().Add(testDefaultTimeout)); err != nil {
		return err
	}

	switch command {
	default:
		return fmt.Errorf("unknown command : %s", command)
	case commandExit:
	case commandVersion:
		_, _ = t.conn.Write(testVersionData)
	case commandStatus3:
		_, _ = t.conn.Write(testStatus3Data)
	case commandLoadStats:
		_, _ = t.conn.Write(testLoadStatsData)
	}
	return nil
}
