package socket

import (
	"crypto/tls"
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
	Network: NetworkTCP,
	Address: testServerAddress,
	Timeout: 10 * time.Millisecond,
	TLSConf: nil,
}

var udpConfig = Config{
	Network: NetworkUDP,
	Address: testServerAddress,
	Timeout: 200 * time.Millisecond,
	TLSConf: nil,
}

var unixConfig = Config{
	Network: NetworkUnix,
	Address: testUnixServerAddress,
	Timeout: 2000 * time.Millisecond,
	TLSConf: nil,
}

var tcpTlsConfig = Config{
	Network: NetworkTCP,
	Address: testServerAddress,
	Timeout: 10 * time.Millisecond,
	TLSConf: &tls.Config{},
}

func Test_clientCommand(t *testing.T) {
	srv := &tcpServer{addr: testServerAddress, rowsNumResp: 1}
	go func() { _ = srv.Run(); defer srv.Close() }()

	time.Sleep(time.Millisecond * 100)
	sock := NewSocket(tcpConfig)
	require.NoError(t, sock.Connect())
	err := sock.Command("ping\n", func(bytes []byte) bool {
		assert.Equal(t, "pong", string(bytes))
		return true
	})
	require.NoError(t, sock.Disconnect())
	require.NoError(t, err)
}

func Test_clientTimeout(t *testing.T) {
	srv := &tcpServer{addr: testServerAddress, rowsNumResp: 1}
	go func() { _ = srv.Run() }()

	time.Sleep(time.Millisecond * 100)
	sock := NewSocket(tcpConfig)
	require.NoError(t, sock.Connect())
	sock.Timeout = 0
	err := sock.Command("ping\n", func(bytes []byte) bool {
		assert.Equal(t, "pong", string(bytes))
		return true
	})
	require.Error(t, err)
}

func Test_clientIncompleteSSL(t *testing.T) {
	srv := &tcpServer{addr: testServerAddress, rowsNumResp: 1}
	go func() { _ = srv.Run() }()

	time.Sleep(time.Millisecond * 100)
	sock := NewSocket(tcpTlsConfig)
	err := sock.Connect()
	require.Error(t, err)
}

func Test_clientCommandStopProcessing(t *testing.T) {
	srv := &tcpServer{addr: testServerAddress, rowsNumResp: 2}
	go func() { _ = srv.Run() }()

	time.Sleep(time.Millisecond * 100)
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
	go func() { _ = srv.Run(); defer srv.Close() }()

	time.Sleep(time.Millisecond * 100)
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
	// cleanup previous file descriptors
	_ = srv.Close()
	go func() { _ = srv.Run() }()

	time.Sleep(time.Millisecond * 200)
	sock := NewSocket(unixConfig)
	require.NoError(t, sock.Connect())
	err := sock.Command("ping\n", func(bytes []byte) bool {
		assert.Equal(t, "pong", string(bytes))
		return false
	})
	require.NoError(t, err)
	require.NoError(t, sock.Disconnect())
}
