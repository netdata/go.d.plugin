package zookeeper

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const maxLinesToRead = 500

func newClient(config clientConfig) *client {
	return &client{network: "tcp", clientConfig: config, dial: net.DialTimeout}
}

type clientConfig struct {
	Address        string
	ReuseRecord    bool
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type dialFunc func(network, address string, timeout time.Duration) (net.Conn, error)

type client struct {
	clientConfig
	dial    dialFunc
	network string
	record  []string
	conn    net.Conn
}

func (c *client) isConnected() bool { return c.conn != nil }

func (c *client) connect() (err error) {
	c.conn, err = c.dial(c.network, c.Address, c.ConnectTimeout)
	return err
}

func (c *client) send(command string) error {
	err := c.conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	if err != nil {
		return err
	}
	_, err = c.conn.Write([]byte(command))
	return err
}

func (c *client) read() (record []string, err error) {
	if err = c.conn.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
		return nil, err
	}

	if c.ReuseRecord {
		record, err = read(c.record, c.conn)
		c.record = record
	} else {
		record, err = read(nil, c.conn)
	}

	return record, err
}

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

		line = strings.Trim(line, "\r\n ")
		dst = append(dst, line)

		num++
		if num > maxLinesToRead {
			err = fmt.Errorf("read line limit exceeded (%d)", maxLinesToRead)
			break
		}
	}
	return dst, err
}
