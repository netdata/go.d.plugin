package chrony

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

func (c *Chrony) SubmitRequest(req *RequestPacket) (*ReplyPacket, *bytes.Reader, error) {
	conn := c.conn
	var err error

	var seqNumber uint32
	if req.SeqNumber != 0 {
		seqNumber = req.SeqNumber
	} else {
		seqNumber = uint32(time.Now().Unix())
		req.SeqNumber = seqNumber
	}

	// request marshal then write
	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		return nil, nil, fmt.Errorf("failed to write request: %s", err)
	}

	// get rsp
	var rspLen int
	dgram := make([]byte, 10240)
	rspLen, err = conn.Read(dgram)
	if err != nil {
		return nil, nil, err
	}
	c.Debugf("read %d byte from response", rspLen)

	rd := bytes.NewReader(dgram)
	var reply ReplyPacket
	if err := binary.Read(rd, binary.BigEndian, &reply); err != nil {
		return nil, nil, fmt.Errorf("failed to get relay from conn: %s", err)
	}
	c.Debugf("req: %+v rsp:%+v\n", req, reply)

	// check every fields
	if reply.SeqNum != seqNumber {
		return &reply, rd, fmt.Errorf("unexpected tracking packet seqNumber: %d", reply.SeqNum)
	}

	if reply.Version != req.Version {
		return &reply, rd, fmt.Errorf("unexpected chrony protocol version: %d", reply.Version)
	}

	return &reply, rd, nil
}

func (c *Chrony) ParseChronyReply(reply *ReplyPacket, rd *bytes.Reader, err error) (*ReplyPacket, interface{}, error) {
	switch reply.PktType {
	case pktTypeCMDReply:
	default:
		return reply, nil, fmt.Errorf("unexpected chrony reply type: %d", reply.PktType)
	}

	// get command from relay then apply
	var payload interface{}
	switch reply.Command {
	case reqActivity:
		payload = &ActivityPayload{}
	case reqTracking:
		payload = &TrackingPayload{}
	default:
		payload = make([]byte, rd.Len())
		err = fmt.Errorf("unexpected reply command: %d", reply.Command)
	}

	// get rsp body
	if err := binary.Read(rd, binary.BigEndian, payload); err != nil {
		return reply, nil, fmt.Errorf("failed reading payload: %s", err)
	}

	return reply, payload, err
}

func (c *Chrony) FetchTracking() (*TrackingPayload, error) {
	req := c.EmptyRequest()
	req.Command = reqTracking

	_, trackingPtr, err := c.ParseChronyReply(c.SubmitRequest(req))
	if err != nil {
		return nil, err
	}

	return trackingPtr.(*TrackingPayload), nil
}

func (c *Chrony) FetchActivity() (*ActivityPayload, error) {
	req := c.EmptyRequest()
	req.Command = reqActivity

	_, activityPtr, err := c.ParseChronyReply(c.SubmitRequest(req))
	if err != nil {
		return nil, err
	}

	return activityPtr.(*ActivityPayload), nil
}

func (c *Chrony) EmptyRequest() *RequestPacket {
	// Check() func would init the value.
	if c.chronyVersion == 0 {
		err := c.ApplyChronyVersion()
		if err != nil {
			panic(err) // unexpected chrony protocol version, we can't collect data correct.
		}
	}
	return &RequestPacket{
		Version: c.chronyVersion,
		PktType: pktTypeCMDRequest,
	}
}

func (c *Chrony) SubmitEmptyRequest() error {
	_, _, err := c.SubmitRequest(c.EmptyRequest())
	return err
}

func (c *Chrony) ApplyChronyVersion() error {

	tryProtocolVersion := []uint8{
		protoVersionNumber6,
		protoVersionNumber5,
	}
	for _, version := range tryProtocolVersion {
		rpy, _, err := c.SubmitRequest(&RequestPacket{
			Version: version,
			PktType: pktTypeCMDRequest,
			Command: 0,
		})
		if err != nil {
			c.Debugf("contact chrony failed with err: %+v", err)
			continue
		}

		c.Debugf("chrony reply protocol version: %d", rpy.Version)
		if version == rpy.Version {
			c.chronyVersion = version
			return nil
		}
	}

	c.Warningf("will use default chrony protocol version")
	c.chronyVersion = protoVersionNumber
	return nil
	//return fmt.Errorf("unexpected chrony protocol version")
}
