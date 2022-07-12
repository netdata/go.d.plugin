package client

import (
	"fmt"
	"math"
	"net"
	"strings"
	"time"
)

const (
	// // https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/candm.h#L375-L376
	// pktTypeCMDRequest is the request packet type
	pktTypeCMDRequest = uint8(1)
	// pktTypeCMDReply is the reply packet type
	pktTypeCMDReply = uint8(2)
)

const (
	// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/candm.h#L39
	// reqTracking identifies a tracking request (REQ_TRACKING)
	reqTracking = uint16(33)
	// reqActivity identifies an activity check request (REQ_ACTIVITY)
	reqActivity = uint16(44)
)

// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/candm.h#L431
// requestHead represents CMD_Request
type requestHead struct {
	Version   uint8
	PktType   uint8
	Res1      uint8
	Res2      uint8
	Command   uint16
	Attempt   uint16
	SeqNumber uint32
	Pad       [396]byte
}

// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/candm.h#L784
// replyHead represents CMD_Reply
type replyHead struct {
	Version uint8
	PktType uint8
	Res1    uint8
	Res2    uint8
	Command uint16
	Reply   uint16
	Status  uint16
	Pad1    uint16
	Pad2    uint16
	Pad3    uint16
	SeqNum  uint32
	Pad4    uint32
	Pad5    uint32
}

// TrackingPayload is the payload for tracking replies (RPY_Tracking)
// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/candm.h#L581
type TrackingPayload struct {
	RefID              uint32
	Ip                 IPAddr
	Stratum            uint16
	LeapStatus         uint16
	RefTime            ChronyTimespec
	CurrentCorrection  ChronyFloat
	LastOffset         ChronyFloat
	RmsOffset          ChronyFloat
	FreqPpm            ChronyFloat
	ResidFreqPpm       ChronyFloat
	SkewPpm            ChronyFloat
	RootDelay          ChronyFloat
	RootDispersion     ChronyFloat
	LastUpdateInterval ChronyFloat
}

func (tp *TrackingPayload) String() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("RefID: %d\n", tp.RefID))
	b.WriteString(fmt.Sprintf("ActiveServer: %s\n", tp.Ip))
	b.WriteString(fmt.Sprintf("Stratum: %d\n", tp.Stratum))
	b.WriteString(fmt.Sprintf("RefTime: %s\n", tp.RefTime.Time().Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("CurrentCorrection: %f\n", tp.CurrentCorrection.Float64()))
	b.WriteString(fmt.Sprintf("FreqPpm: %f\n", tp.FreqPpm.Float64()))
	b.WriteString(fmt.Sprintf("ResidFreqPpm: %f\n", tp.ResidFreqPpm.Float64()))
	b.WriteString(fmt.Sprintf("SkewPpm: %f\n", tp.SkewPpm.Float64()))
	b.WriteString(fmt.Sprintf("RootDelay: %f\n", tp.RootDelay.Float64()))
	b.WriteString(fmt.Sprintf("RootDispersion: %f\n", tp.RootDispersion.Float64()))
	b.WriteString(fmt.Sprintf("LeapStatus: %d\n", tp.LeapStatus))
	b.WriteString(fmt.Sprintf("LastUpdateInterval: %f\n", tp.LastUpdateInterval.Float64()))
	b.WriteString(fmt.Sprintf("LastOffset: %f\n", tp.LastOffset.Float64()))
	b.WriteString(fmt.Sprintf("RmsOffset: %f", tp.RmsOffset.Float64()))
	return b.String()
}

// ActivityPayload is the payload for activity replies (RPY_Activity)
// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/candm.h#L685
type ActivityPayload struct {
	Online       int32
	Offline      int32
	BurstOnline  int32
	BurstOffline int32
	Unresolved   int32
}

// ChronyTimespec is the custom chrony timespec type (Timespec)
// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/candm.h#L115
type ChronyTimespec struct {
	TvSecHigh uint32
	TvSecLow  uint32
	TvNSec    uint32
}

func (ct ChronyTimespec) Time() time.Time {
	nsec := uint32(999999999)
	if ct.TvNSec < nsec {
		nsec = ct.TvNSec
	}
	return time.Unix(int64(uint64(ct.TvSecHigh)<<32+uint64(ct.TvSecLow)), int64(nsec))
}

const (
	// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/util.c#L891
	// floatExpBits represents 32-bit floating-point format consisting of 7-bit signed exponent
	floatExpBits = 7

	// floatCoefBits represents chronyFloat 25-bit signed coefficient without hidden bit
	floatCoefBits = 25

	ScaleFactor = 1000000000
)

// ChronyFloat is 32-bit floating-point format consisting of 7-bit signed exponent
// and 25-bit signed coefficient without hidden bit.
// The result is calculated as: 2^(exp - 25) * coef.
type ChronyFloat int32

// Float64 does magic to decode float from int32.
// https://github.com/mlichvar/chrony/blob/2ac22477563581ae3bc39c4ff28464059c0a73be/util.c#L900
func (cf ChronyFloat) Float64() float64 {
	var exp, coef int32

	x := uint32(cf)

	exp = int32(x >> floatCoefBits)
	if exp >= 1<<(floatExpBits-1) {
		exp -= 1 << floatExpBits
	}
	exp -= floatCoefBits

	coef = int32(x % (1 << floatCoefBits))
	if coef >= 1<<(floatCoefBits-1) {
		coef -= 1 << floatCoefBits
	}

	return float64(coef) * math.Pow(2.0, float64(exp))
}

// Int64 returns the 64bits float value
func (cf ChronyFloat) Int64() int64 { return int64(cf.Float64() * ScaleFactor) }

// IPAddr represents IPAddr structure.
// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/addressing.h#L41
type IPAddr struct {
	IPAddrHigh uint64
	IPAddrLow  uint64
	Family     uint16
	Pad        uint16
}

func (ia IPAddr) String() string { return ia.IP().String() }

func (ia IPAddr) IP() net.IP {
	const ipAddrInet4 = uint16(1)
	const ipAddrInet6 = uint16(2)

	if ia.Family == ipAddrInet4 {
		m := uint32(ia.IPAddrHigh >> (32))
		var ip [4]uint8
		for i := 0; i < 4; i++ {
			ip[i] = uint8(m % 0x100)
			m = m / 0x100
		}
		return net.IPv4(ip[3], ip[2], ip[1], ip[0])
	}

	if ia.Family == ipAddrInet6 {
		addr := make(net.IP, net.IPv6len)
		h := ia.IPAddrHigh
		for i := 7; i >= 0; i-- {
			addr[i] = byte(h % 0x100)
			h = h / 0x100
		}
		l := ia.IPAddrLow
		for i := 7; i >= 0; i-- {
			addr[i+8] = byte(l % 0x100)
			l = l / 0x100
		}
		return addr
	}

	return net.IPv4zero
}
