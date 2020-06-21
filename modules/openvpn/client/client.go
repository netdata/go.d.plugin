package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reLoadStats = regexp.MustCompile(`^SUCCESS: nclients=([0-9]+),bytesin=([0-9]+),bytesout=([0-9]+)`)
	reVersion   = regexp.MustCompile(`^OpenVPN Version: OpenVPN ([0-9]+)\.([0-9]+)\.([0-9]+) .+Management Version: ([0-9])`)
)

const maxLinesToRead = 500

// New creates new OpenVPN client.
func New(config Config) *Client {
	network := "tcp"
	if strings.HasPrefix(config.Address, "/") {
		network = "unix"
	}
	return &Client{network: network, Config: config, dial: net.DialTimeout}
}

// Config represents OpenVPN client config.
type Config struct {
	Address        string
	ReuseRecord    bool
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type dialFunc func(network, address string, timeout time.Duration) (net.Conn, error)

// Client represents OpenVPN client.
type Client struct {
	Config
	dial    dialFunc
	network string
	record  []string
	conn    net.Conn
}

// IsConnected IsConnected.
func (c *Client) IsConnected() bool { return c.conn != nil }

// Connect connects.
func (c *Client) Connect() (err error) {
	if c.IsConnected() {
		_ = c.Disconnect()
	}
	c.conn, err = c.dial(c.network, c.Address, c.ConnectTimeout)
	return err
}

// Disconnect closes connection, if there is no connection it does nothing.
func (c *Client) Disconnect() error {
	if !c.IsConnected() {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

// Users Users.
func (c *Client) Users() (Users, error) {
	lines, err := c.get(commandStatus3, readUntilEND)
	if err != nil {
		return nil, err
	}
	return decodeUsers(lines)
}

// LoadStats LoadStats.
func (c *Client) LoadStats() (*LoadStats, error) {
	lines, err := c.get(commandLoadStats, readOneLine)
	if err != nil {
		return nil, err
	}
	return decodeLoadStats(lines)
}

// Version Version.
func (c *Client) Version() (*Version, error) {
	lines, err := c.get(commandVersion, readUntilEND)
	if err != nil {
		return nil, err
	}
	return decodeVersion(lines)
}

func (c *Client) get(command string, stopRead stopReadFunc) ([]string, error) {
	if err := c.send(command); err != nil {
		return nil, err
	}
	return c.read(stopRead)
}

func (c *Client) send(command string) error {
	if !c.IsConnected() {
		return errors.New("not connected")
	}
	err := c.conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	if err != nil {
		return err
	}
	_, err = c.conn.Write([]byte(command))
	return err
}

func (c *Client) read(stopRead stopReadFunc) (record []string, err error) {
	if !c.IsConnected() {
		return nil, errors.New("not connected")
	}

	if err = c.conn.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
		return nil, err
	}

	if c.ReuseRecord {
		record, err = read(c.record, c.conn, stopRead)
		c.record = record
	} else {
		record, err = read(nil, c.conn, stopRead)
	}

	return record, err
}

func read(dst []string, reader io.Reader, stopRead stopReadFunc) ([]string, error) {
	dst = dst[:0]
	var (
		r    = bufio.NewReader(reader)
		err  error
		num  int
		line string
	)
	for {
		line, err = r.ReadString('\n')
		// never EOF
		if err != nil {
			break
		}
		// skip real-time messages
		if strings.HasPrefix(line, ">") {
			continue
		}

		line = strings.Trim(line, "\r\n ")
		dst = append(dst, line)
		if stopRead != nil && stopRead(line) {
			break
		}

		num++
		if num > maxLinesToRead {
			err = fmt.Errorf("read line limit exceeded (%d)", maxLinesToRead)
			break
		}
	}
	return dst, err
}

type stopReadFunc func(string) bool

func readOneLine(_ string) bool { return true }

func readUntilEND(s string) bool { return strings.HasSuffix(s, "END") }

func decodeLoadStats(src []string) (*LoadStats, error) {
	m := reLoadStats.FindStringSubmatch(strings.Join(src, " "))
	if len(m) == 0 {
		return nil, fmt.Errorf("parse failed : %v", src)
	}
	return &LoadStats{
		NumOfClients: mustParseInt(m[1]),
		BytesIn:      mustParseInt(m[2]),
		BytesOut:     mustParseInt(m[3]),
	}, nil
}

func decodeVersion(src []string) (*Version, error) {
	m := reVersion.FindStringSubmatch(strings.Join(src, " "))
	if len(m) == 0 {
		return nil, fmt.Errorf("parse failed : %v", src)
	}
	return &Version{
		Major:      mustParseInt(m[1]),
		Minor:      mustParseInt(m[2]),
		Patch:      mustParseInt(m[3]),
		Management: mustParseInt(m[4]),
	}, nil
}

// works only for `status 3\n`
func decodeUsers(src []string) (Users, error) {
	var users Users

	// [CLIENT_LIST common_name 178.66.34.194:54200 10.9.0.5 9319 8978 Thu May 9 05:01:44 2019 1557345704 username]
	for _, v := range src {
		if !strings.HasPrefix(v, "CLIENT_LIST") {
			continue
		}
		parts := strings.Fields(v)
		// Right after the connection there are no virtual ip, and both common name and username UNDEF
		// CLIENT_LIST	UNDEF	178.70.95.93:39324		1411	3474	Fri May 10 07:41:54 2019	1557441714	UNDEF
		if len(parts) != 13 {
			continue
		}
		u := User{
			CommonName:     parts[1],
			RealAddress:    parts[2],
			VirtualAddress: parts[3],
			BytesReceived:  mustParseInt(parts[4]),
			BytesSent:      mustParseInt(parts[5]),
			ConnectedSince: mustParseInt(parts[11]),
			Username:       parts[12],
		}
		users = append(users, u)
	}
	return users, nil
}

func mustParseInt(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}
