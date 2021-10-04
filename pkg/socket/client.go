package socket

import (
	"bufio"
	"crypto/tls"
	"errors"
	"net"
	"time"
)

func NewSocket(config Config) *socket {
	return &socket{
		Config: config,
		conn:   nil,
	}
}

type socket struct {
	Config
	conn net.Conn
}

func (s *socket) Connect() (err error) {
	if s.TLSConf == nil {
		s.conn, err = net.DialTimeout(string(s.Network), s.Address, s.Timeout)
	} else {
		var d net.Dialer
		d.Timeout = s.Timeout
		s.conn, err = tls.DialWithDialer(&d, string(s.Network), s.Address, s.TLSConf)
	}
	return err
}

func (s *socket) Disconnect() (err error) {
	if s.conn != nil {
		err = s.conn.Close()
		s.conn = nil
	}
	return err
}

func (s *socket) Command(command string, process Processor) error {
	if err := write(command, s.conn, s.Timeout); err != nil {
		return err
	}
	return read(s.conn, process, s.Timeout)
}

func write(command string, writer net.Conn, timeout time.Duration) error {
	if err := writer.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	_, err := writer.Write([]byte(command))
	return err
}

func read(reader net.Conn, process Processor, timeout time.Duration) error {
	if process == nil {
		return errors.New("process func is nil")
	}
	if err := reader.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() && process(scanner.Bytes()) {
	}
	return scanner.Err()
}
