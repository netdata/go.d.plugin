package socket

import (
	"bufio"
	"crypto/tls"
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
	if s.TlsConf == nil {
		s.conn, err = net.DialTimeout(string(s.Network), s.Address, s.Timeout)
	} else {
		var d net.Dialer
		d.Timeout = s.Timeout
		s.conn, err = tls.DialWithDialer(&d, string(s.Network), s.Address, s.TlsConf)
	}
	return err
}

func (s *socket) Disconnect() error {
	err := s.conn.Close()
	s.conn = nil
	return err
}

func (s *socket) Command(command string, process Processor) error {
	if err := s.send(command, s.conn, s.Timeout); err != nil {
		return err
	}
	return read(s.conn, process, s.Timeout)
}

func (s *socket) send(command string, writer net.Conn, timeout time.Duration) error {
	err := writer.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(command))
	return err
}

func read(reader net.Conn, process Processor, timeout time.Duration) error {
	err := reader.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	s := bufio.NewScanner(reader)
	for s.Scan() && process(s.Bytes()) {
	}
	return s.Err()
}
