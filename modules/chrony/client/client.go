package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/netdata/go.d.plugin/logger"
)

const (
	// protoVersionNumber is the protocol version for this client
	protoVersionNumber  = protoVersionNumber6
	protoVersionNumber6 = uint8(6)
	protoVersionNumber5 = uint8(5)
)

type Config struct {
	Address string
	Timeout time.Duration
}

func New(l *logger.Logger, c Config) (*Client, error) {
	conn, err := net.DialTimeout("udp", c.Address, c.Timeout)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Logger:  l,
		conn:    conn,
		timeout: c.Timeout,
	}
	client.chronyVersion = client.guessChronyVersion()

	return client, nil
}

type Client struct {
	*logger.Logger
	conn          net.Conn
	timeout       time.Duration
	chronyVersion uint8
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
}

func (c *Client) Tracking() (*TrackingPayload, error) {
	req := &requestHead{
		Version: c.chronyVersion,
		PktType: pktTypeCMDRequest,
		Command: reqTracking,
	}

	_, payload, err := c.query(req)
	if err != nil {
		return nil, err
	}

	var tp TrackingPayload
	if err := binary.Read(payload, binary.BigEndian, &tp); err != nil {
		return nil, fmt.Errorf("failed reading tracking payload: %s", err)
	}

	return &tp, nil
}

func (c *Client) Activity() (*ActivityPayload, error) {
	req := &requestHead{
		Version: c.chronyVersion,
		PktType: pktTypeCMDRequest,
		Command: reqActivity,
	}

	_, payload, err := c.query(req)
	if err != nil {
		return nil, err
	}

	var ap ActivityPayload
	if err := binary.Read(payload, binary.BigEndian, &ap); err != nil {
		return nil, fmt.Errorf("failed reading activity reply: %s", err)
	}

	return &ap, nil
}

func (c *Client) Ping() error {
	req := &requestHead{
		Version: c.chronyVersion,
		PktType: pktTypeCMDRequest,
	}
	_, _, err := c.query(req)
	return err
}

func (c *Client) guessChronyVersion() uint8 {
	versions := []uint8{
		protoVersionNumber6,
		protoVersionNumber5,
	}
	for _, ver := range versions {
		req := &requestHead{
			Version: ver,
			PktType: pktTypeCMDRequest,
		}
		rpy, _, err := c.query(req)
		if err != nil {
			c.Debugf("contact chrony failed with err: %+v", err)
			continue
		}
		if ver == rpy.Version {
			c.Debugf("chrony reply protocol version: %d", ver)
			return ver
		}
	}

	c.Warningf("will use default chrony protocol version: %d", protoVersionNumber)
	return protoVersionNumber
}

func (c *Client) query(req *requestHead) (*replyHead, *bytes.Reader, error) {
	if req.Version == 0 {
		return nil, nil, errors.New("request version is not set")
	}

	if req.SeqNumber == 0 {
		req.SeqNumber = uint32(time.Now().Unix())
	}

	if err := c.conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		return nil, nil, err
	}

	if err := binary.Write(c.conn, binary.BigEndian, req); err != nil {
		return nil, nil, fmt.Errorf("failed to write request: %v", err)
	}

	if err := c.conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		return nil, nil, err
	}

	dgram := make([]byte, 1024)
	if _, err := c.conn.Read(dgram); err != nil {
		return nil, nil, err
	}

	payload := bytes.NewReader(dgram)
	var rpy replyHead
	if err := binary.Read(payload, binary.BigEndian, &rpy); err != nil {
		return nil, nil, fmt.Errorf("failed to read rpy: %v", err)
	}

	if rpy.PktType != pktTypeCMDReply {
		return &rpy, payload, fmt.Errorf("unexpected packet type: want=%d, got=%d", pktTypeCMDReply, rpy.PktType)
	}
	if rpy.SeqNum != req.SeqNumber {
		return &rpy, payload, fmt.Errorf("unexpected rpy seqNumber: want=%d, got=%d", req.SeqNumber, rpy.SeqNum)
	}
	if rpy.Version != req.Version {
		return &rpy, payload, fmt.Errorf("unexpected rpy protocol version: want=%d, got=%d", req.Version, rpy.Version)
	}

	return &rpy, payload, nil
}
