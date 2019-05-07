package openvpn

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testLoadStatsData, _    = ioutil.ReadFile("testdata/load-stats.txt")
	testVersionData, _      = ioutil.ReadFile("testdata/version.txt")
	testStatus3Data, _      = ioutil.ReadFile("testdata/status3.txt")
	testStatus3EmptyData, _ = ioutil.ReadFile("testdata/status3-empty.txt")

	testCommandStatus3Empty = "status 3 empty/n"
)

func TestNew(t *testing.T) {
	job := New()
	assert.IsType(t, (*OpenVPN)(nil), job)
	assert.Equal(t, defaultAddress, job.Address)
	assert.Equal(t, defaultConnectTimeout, job.Timeouts.Connect.Duration)
	assert.Equal(t, defaultReadTimeout, job.Timeouts.Read.Duration)
	assert.Equal(t, defaultWriteTimeout, job.Timeouts.Write.Duration)
}

func TestOpenVPN_Init(t *testing.T) {
	job := New()
	assert.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestOpenVPN_Check(t *testing.T) {
	job := New()
	assert.True(t, job.Init())
	serverConn, clientConn := net.Pipe()
	job.apiClient = newTestClient(clientConn)
	server := newTestServer(serverConn)
	go server.serve()
	assert.True(t, job.Check())
}

func TestOpenVPN_CheckNG(t *testing.T) {
	job := New()
	job.Address = "127.0.0.1:38001"
	assert.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestOpenVPN_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestOpenVPN_Cleanup(t *testing.T) { assert.NotPanics(t, New().Cleanup) }

func TestExample_Collect(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
	serverConn, clientConn := net.Pipe()
	job.apiClient = newTestClient(clientConn)
	server := newTestServer(serverConn)
	go server.serve()

	assert.True(t, job.Check())
	expected := map[string]int64{
		"bytes_in":  7811,
		"bytes_out": 7667,
		"clients":   1,
	}
	assert.Equal(t, expected, job.Collect())
}

func newTestClient(conn net.Conn) *client {
	return &client{
		conn: conn,
		clientConfig: clientConfig{
			timeouts: clientTimeouts{
				connect: defaultConnectTimeout,
				read:    defaultReadTimeout,
				write:   defaultWriteTimeout,
			},
		},
	}
}

func newTestServer(conn net.Conn) *testTCPServer { return &testTCPServer{conn: conn} }

type testTCPServer struct {
	conn net.Conn
}

func (t *testTCPServer) serve() error {
	for {
		err := t.serveOnce()
		if err != nil {
			return err
		}
	}
}

func (t *testTCPServer) serveOnce() error {
	_ = t.conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
	s, err := bufio.NewReader(t.conn).ReadString('\n')
	if err != nil {
		return err
	}
	err = t.conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
	if err != nil {
		return err
	}
	switch s {
	default:
		return fmt.Errorf("unknown command : %s", s)
	case commandExit:
		return errors.New("exiting")
	case commandVersion:
		_, _ = t.conn.Write(testVersionData)
	case commandStatus:
		_, _ = t.conn.Write(testStatus3Data)
	case commandLoadStats:
		_, _ = t.conn.Write(testLoadStatsData)
	case testCommandStatus3Empty:
		_, _ = t.conn.Write(testStatus3EmptyData)
	}
	return nil
}
