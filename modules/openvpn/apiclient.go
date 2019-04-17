package openvpn

import (
	"bufio"
	"fmt"
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
	commandStatus    = "status 3\n"   // --status-version 3
	commandLoadStats = "load-stats\n" // no description in docs ¯\(°_o)/¯
	// Show the current OpenVPN and Management Interface versions.
	commandVersion = "version\n"
)

type apiClient interface {
	connect() error
	reconnect() error
	disconnect() error
	send(command string) error
	read(stop func(string) bool) ([]string, error)
	isConnected() bool
}

func newClient(config clientConfig) *client {
	return &client{clientConfig: config}
}

type clientConfig struct {
	network  string
	address  string
	timeouts clientTimeouts
}

type clientTimeouts struct {
	connect time.Duration
	read    time.Duration
	write   time.Duration
}

type client struct {
	clientConfig
	records []string
	conn    net.Conn
}

func (c *client) isConnected() bool {
	return c.conn != nil
}

func (c *client) connect() error {
	if c.conn != nil {
		return c.reconnect()
	}
	conn, err := net.DialTimeout(c.network, c.address, c.timeouts.connect)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) reconnect() error {
	if c.conn != nil {
		_ = c.disconnect()
	}
	return c.connect()
}

func (c *client) disconnect() error {
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *client) send(command string) error {
	fmt.Println("SEND", command)
	err := c.conn.SetWriteDeadline(time.Now().Add(c.timeouts.write))
	if err != nil {
		return err
	}
	_, err = c.conn.Write([]byte(command))
	return err
}

func (c *client) read(stop func(string) bool) ([]string, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(c.timeouts.read))
	if err != nil {
		return nil, err
	}
	c.records = c.records[:0]
	r := bufio.NewReader(c.conn)
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
		line = strings.Trim(line, "\r\n ")
		c.records = append(c.records, line)
		if stop != nil && stop(line) {
			break
		}
	}
	return c.records, nil
}
