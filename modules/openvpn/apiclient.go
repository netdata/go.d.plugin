package openvpn

import (
	"bufio"
	"net"
	"strings"
	"time"
)

/*
https://openvpn.net/community-resources/management-interface/

OUTPUT FORMAT
-------------

(1) Command success/failure indicated by "SUCCESS: [text]" or
    "ERROR: [text]".

(2) For commands which print multiple lines of output,
    the last line will be "END".

(3) Real-time messages will be in the form ">[source]:[text]",
    where source is "CLIENT", "ECHO", "FATAL", "HOLD", "INFO", "LOG",
    "NEED-OK", "PASSWORD", or "STATE".
*/

var (
	// Close the management session, and resume listening on the
	// management port for connections from other clients. Currently,
	// the OpenVPN daemon can at most support a single management client
	// any one time.
	commandExit = "exit\n" // "quit"
	// Show current daemon status information, in the same format as
	// that produced by the OpenVPN --status directive.
	commandStatus    = "status 3\n" // --status-version 3
	commandLoadStats = "load-stats\n"
)

func newAPIClient(config apiClientConfig) *apiClient {
	return &apiClient{apiClientConfig: config}
}

type apiClientConfig struct {
	network string
	address string
	timeout struct {
		connect, read time.Duration
	}
}

type apiClient struct {
	apiClientConfig

	resp []string
	conn net.Conn
}

func (a *apiClient) connect() error {
	if a.conn != nil {
		return a.reconnect()
	}
	conn, err := net.DialTimeout(a.network, a.address, a.timeout.connect)
	if err != nil {
		return err
	}
	a.conn = conn
	return nil
}

func (a *apiClient) reconnect() error {
	if a.conn != nil {
		_ = a.disconnect()
	}
	return a.connect()
}

func (a *apiClient) disconnect() error {
	if a.conn == nil {
		return nil
	}
	_ = a.send(commandExit)
	err := a.conn.Close()
	a.conn = nil
	return err
}

func (a *apiClient) send(command string) error {
	_, err := a.conn.Write([]byte(command))
	return err
}

func (a *apiClient) read(stop func(string) bool) ([]string, error) {
	err := a.conn.SetReadDeadline(time.Now().Add(a.timeout.read))
	if err != nil {
		return nil, err
	}
	a.resp = a.resp[:0]
	r := bufio.NewReader(a.conn)
	var line string
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			break
		}
		// skip real-time messages
		if strings.HasPrefix(line, ">") {
			continue
		}
		a.resp = append(a.resp, line)
		if stop(line) {
			break
		}
	}
	return a.resp, nil
}
