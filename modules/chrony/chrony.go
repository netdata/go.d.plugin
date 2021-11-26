package chrony

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	// "github.com/getsentry/raven-go"

	"github.com/getsentry/sentry-go"
	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	// protoVersionNumber is the protocol version for this client
	protoVersionNumber = uint8(6)

	// pktTypeCMDRequest is the request packet type
	pktTypeCMDRequest = uint8(1)
	// pktTypeCMDReply is the reply packet type
	pktTypeCMDReply = uint8(2)

	// reqTracking identifies a tracking request
	reqTracking = uint16(33)
	// reqActivity identifies an activity check request
	reqActivity = uint16(44)
	// reqNSources identifies a n_sources request
	reqNSources = uint16(14)
	// reqSourceStats identifies a sourcestats request
	reqSourceStats = uint16(34)
	// reqSourceData identifies a source data request
	reqSourceData = uint16(15)

	// rpyTracking identifies a tracking reply
	rpyTypeTracking = uint16(5)
	// rpyActivity identifies an activity check reply
	rpyTypeActivity = uint16(12)
	// reqNSources identifies a n_sources request
	rpyNSources = uint16(2)
	// reqSourcesStats identifies a sourcestats request
	rpySourcesStats = uint16(6)

	// floatExpBits represent 32-bit floating-point format consisting of 7-bit signed exponent
	floatExpBits = 7
	// floatCoefBits represent chronyFloat 25-bit signed coefficient without hidden bit
	floatCoefBits = 25
	// precision scaling factor
	scaleFactor = 1000000000
)

var (
	// chronyCmdAddr is the chrony local port
	chronyCmdAddr = "127.0.0.1:323"
)

func init() {
	// raven.SetDSN("https://a646445445cc43168fd66a4095c17526:9df6ff045c7440f196900e573d13903c@sentry.i.agoralab.co/23")
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("chrony", creator)
}

// New creates chronyCollector exposing local status of a chrony daemon
func New() *chronyCollector {
	addr, _ := net.ResolveUDPAddr("udp", chronyCmdAddr)
	return &chronyCollector{
		metrics: make(map[string]int64),
		cmdAddr: addr,
	}
}

// Cleanup makes cleanup
func (c *chronyCollector) Cleanup() {
}

// Init makes initialization
func (c *chronyCollector) Init() bool {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://a646445445cc43168fd66a4095c17526:9df6ff045c7440f196900e573d13903c@sentry.i.agoralab.co/23",
	})

	if err != nil {
		c.Errorf("Sentry initialization failed: %v\n", err)
	}
	// sentry.Flush(time.Second * 5)
	return true
}

// Check makes check
func (c *chronyCollector) Check() bool {
	if c.cmdAddr == nil {
		c.Errorf("invalid chrony cmdAddr %s", chronyCmdAddr)
		return false
	}
	// TODO need something real to check udp socket
	conn, err := net.DialUDP("udp", nil, c.cmdAddr)
	if err != nil {
		c.Errorf("unable connect to chrony addr %s, is chrony up and running?", c.cmdAddr)
		return false
	}

	// TODO: need protocol version check
	defer conn.Close()
	c.Debugf("connecting to chrony addr %s", conn.LocalAddr())
	return true
}

// Charts creates Charts dynamically (each gpu device will be a family)
func (c *chronyCollector) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (c *chronyCollector) Collect() map[string]int64 {
	// collect all we need and sent Exception to sentry
	c.collectTracking()
	c.collectActivity()
	return c.metrics
}

func (c *chronyCollector) collectTracking() {
	tracking, err := c.FetchTracking()
	if err != nil {
		c.Errorf("fetch tracking status failed: %s", err)
		sentry.CaptureException((FetchingChronyError)(err.Error()))
		c.metrics["running"] = 0
		return
	}
	c.Debug(tracking.String())

	c.metrics["running"] = 1
	c.metrics["stratum"] = (int64)(tracking.Stratum)
	c.metrics["leap_status"] = (int64)(tracking.LeapStatus)
	c.metrics["root_delay"] = (int64)(tracking.RootDelay.Int64())
	c.metrics["root_dispersion"] = (int64)(tracking.RootDispersion.Int64())
	c.metrics["skew"] = (int64)(tracking.SkewPpm.Int64())
	c.metrics["frequency"] = (int64)(tracking.FreqPpm.Int64())
	c.metrics["last_offset"] = (int64)(tracking.LastOffset.Int64())
	c.metrics["rms_offset"] = (int64)(tracking.RmsOffset.Int64())
	c.metrics["update_interval"] = (int64)(tracking.LastUpdateInterval.Int64())
	c.metrics["current_correction"] = (int64)(tracking.LastUpdateInterval.Int64())
	c.metrics["ref_timestamp"] = tracking.RefTime.Time().Unix()

	// report root dispersion error to sentry when error > 100ms
	rd := tracking.RootDispersion.Float64()
	if rd > 0.1 {
		c.Debugf("sending sentry error for RootDispersionTooLargeError: %g", rd)
		// raven.CaptureError((RootDispersionTooLargeError)(rd), map[string]string{"service": "chrony"})
		sentry.CaptureException((RootDispersionTooLargeError)(rd))
	}

	// report frequency change to sentry when step > 500ppm
	fp := tracking.FreqPpm.Float64()
	if fp > 500 {
		c.Debugf("sending sentry error for FreqChangeTooFastError: %g", fp)
		sentry.CaptureException((FreqChangeTooFastError)(fp))
	}

	if tracking.LeapStatus > 0 {
		c.Debugf("sending sentry error for LeapStatusError: %g", tracking.LeapStatus)
		sentry.CaptureException((LeapStatusError)(tracking.LeapStatus))
	}
}

func (c *chronyCollector) collectActivity() {
	activity, err := c.FetchActivity()
	if err != nil {
		c.Errorf("fetch activity status failed: %s", err)
		sentry.CaptureException((FetchingChronyError)(err.Error()))
		return
	}
	c.Debug(activity.String())

	c.metrics["online_sources"] = int64(activity.Online)
	c.metrics["offline_sources"] = int64(activity.Offline)
	c.metrics["burst_online_sources"] = int64(activity.BurstOnline)
	c.metrics["burst_offline_sources"] = int64(activity.BurstOffline)
	c.metrics["unresolved_sources"] = int64(activity.Unresolved)

	if activity.Online == 0 {

	}
}

func (c *chronyCollector) SubmitRequest(req *RequestPacket) (*ReplyPacket, interface{}, error) {
	conn, err := net.DialUDP("udp", nil, c.cmdAddr)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

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

	rd := bytes.NewReader(dgram)
	var reply ReplyPacket
	if err := binary.Read(rd, binary.BigEndian, &reply); err != nil {
		return nil, nil, fmt.Errorf("failed to get relay from conn: %s", err)
	}
	c.Debugf("req: %+v rsp:%+v\n", req, reply)

	// check every fields
	if reply.SeqNum != seqNumber {
		return &reply, nil, fmt.Errorf("unexpected tracking packet seqNumber: %d", reply.SeqNum)
	}

	switch reply.ProtoVer {
	case protoVersionNumber:
	default:
		return &reply, nil, fmt.Errorf("unexpected chrony protocol version: %d", reply.ProtoVer)
	}

	switch reply.PktType {
	case pktTypeCMDReply:
	default:
		return &reply, nil, fmt.Errorf("unexpected chrony protocol version: %d", reply.ProtoVer)
	}

	// get command from relay then apply
	var payload interface{}
	switch reply.Command {
	case reqActivity:
		payload = &ActivityPayload{}
	case reqTracking:
		payload = &TrackingPayload{}
	default:
		payload = make([]byte, rspLen-(int(rd.Size())-rd.Len()))
		err = fmt.Errorf("unexpected reply command: %d", reply.Command)
	}

	// get rsp body
	if err := binary.Read(rd, binary.BigEndian, payload); err != nil {
		return &reply, nil, fmt.Errorf("failed reading payload: %s", err)
	}

	return &reply, payload, err
}

func (c *chronyCollector) FetchTracking() (*TrackingPayload, error) {
	var attempt uint16

	req := RequestPacket{
		Version: protoVersionNumber,
		PktType: pktTypeCMDRequest,
		Command: reqTracking,
		Attempt: attempt,
	}

	_, trackingPtr, err := c.SubmitRequest(&req)
	if err != nil {
		return nil, err
	}

	return trackingPtr.(*TrackingPayload), nil
}

func (c *chronyCollector) FetchActivity() (*ActivityPayload, error) {
	var attempt uint16

	req := RequestPacket{
		Version: protoVersionNumber,
		PktType: pktTypeCMDRequest,
		Command: reqActivity,
		Attempt: attempt,
	}

	_, activityPtr, err := c.SubmitRequest(&req)
	if err != nil {
		return nil, err
	}

	return activityPtr.(*ActivityPayload), nil
}
