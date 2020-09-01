package zookeeper

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
	"unsafe"
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
		tlsConf:     config.tlsConf,
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
	return tls.DialWithDialer(&d, c.network, c.address, c.tlsConf)
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
	defer func() {
		_ = c.disconnect()
	}()

	err = c.send(command)
	if err != nil {
		return nil, err
	}

	return c.read()
}

const limitReadLines = 2000

func read(dst []string, reader io.Reader) ([]string, error) {
	dst = dst[:0]
	var num int
	s := bufio.NewScanner(reader)

	for s.Scan() {
		if !isZKLine(s.Bytes()) || isMntrLineOK(s.Bytes()) {
			dst = append(dst, s.Text())
		}
		if num += 1; num >= limitReadLines {
			return nil, fmt.Errorf("read line limit exceeded (%d)", limitReadLines)
		}
	}
	return dst, s.Err()
}

func isZKLine(line []byte) bool {
	return bytes.HasPrefix(line, []byte("zk_"))
}

func isMntrLineOK(line []byte) bool {
	idx := bytes.LastIndexByte(line, '\t')
	return idx > 0 && collectedZKKeys[unsafeString(line)[:idx]]
}

func unsafeString(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}

var collectedZKKeys = map[string]bool{
	"zk_num_alive_connections":      true,
	"zk_outstanding_requests":       true,
	"zk_min_latency":                true,
	"zk_avg_latency":                true,
	"zk_max_latency":                true,
	"zk_packets_received":           true,
	"zk_packets_sent":               true,
	"zk_open_file_descriptor_count": true,
	"zk_max_file_descriptor_count":  true,
	"zk_znode_count":                true,
	"zk_ephemerals_count":           true,
	"zk_watch_count":                true,
	"zk_approximate_data_size":      true,
	"zk_server_state":               true,
}
