package zookeeper

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
)

type clientConfig struct {
	network string
	address string
	timeout time.Duration
	useTLS  bool
	tlsConf *tls.Config
}

func newClient(config clientConfig) *client {
	return &client{
		network:     config.network,
		address:     config.address,
		timeout:     config.timeout,
		useTLS:      config.useTLS,
		tlsConf:     config.tlsConf.Clone(),
		reuseRecord: false,
		record:      nil,
		conn:        nil,
	}
}

type client struct {
	network     string
	address     string
	timeout     time.Duration
	useTLS      bool
	tlsConf     *tls.Config
	reuseRecord bool
	record      []string
	conn        net.Conn
}

func (c client) dial() (net.Conn, error) {
	if !c.useTLS {
		return net.DialTimeout(c.network, c.address, c.timeout)
	}
	var d net.Dialer
	d.Timeout = c.timeout
	return tls.DialWithDialer(&d, c.network, c.address, c.tlsConf.Clone())
}

func (c *client) connect() (err error) {
	c.conn, err = c.dial()
	return err
}

func (c *client) disconnect() error {
	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *client) isConnected() bool {
	return c.conn != nil
}

func (c *client) send(command string) error {
	err := c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return err
	}
	_, err = c.conn.Write([]byte(command))
	return err
}

func (c *client) read() (record []string, err error) {
	if err = c.conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		return nil, err
	}

	if c.reuseRecord {
		record, err = read(c.record, c.conn)
		c.record = record
	} else {
		record, err = read(nil, c.conn)
	}

	return record, err
}

func (c *client) fetch(command string) (rows []string, err error) {
	if c.isConnected() {
		_ = c.disconnect()
	}

	err = c.connect()
	if err != nil {
		return nil, err
	}
	defer c.disconnect()

	err = c.send(command)
	if err != nil {
		return nil, err
	}

	return c.read()
}

const maxLinesToRead = 5001

func read(dst []string, reader io.Reader) ([]string, error) {
	dst = dst[:0]
	var (
		r    = bufio.NewReader(reader)
		err  error
		num  int
		line string
	)
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		dst = append(dst, line)

		num++
		if num > maxLinesToRead {
			err = fmt.Errorf("read line limit exceeded (%d)", maxLinesToRead)
			break
		}
	}
	return dst, err
}
