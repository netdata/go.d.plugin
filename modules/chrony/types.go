package chrony

import (
	"fmt"
	"math"
	"net"
	"time"
)

const (
	// protoVersionNumber is the protocol version for this client
	protoVersionNumber  = protoVersionNumber6
	protoVersionNumber6 = uint8(6)
	protoVersionNumber5 = uint8(5)

	// pktTypeCMDRequest is the request packet type
	pktTypeCMDRequest = uint8(1)
	// pktTypeCMDReply is the reply packet type
	pktTypeCMDReply = uint8(2)

	// reqTracking identifies a tracking request
	reqTracking = uint16(33)
	// reqActivity identifies an activity check request
	reqActivity = uint16(44)

	// floatExpBits represent 32-bit floating-point format consisting of 7-bit signed exponent
	floatExpBits = 7
	// floatCoefBits represent chronyFloat 25-bit signed coefficient without hidden bit
	floatCoefBits = 25

	scaleFactor = 1000000000
)

// RequestPacket holds a chrony request
type requestPacket struct {
	Version   uint8 /* Protocol version */
	PktType   uint8 /* What sort of packet this is */
	Res1      uint8
	Res2      uint8
	Command   uint16 /* Which command is being issued */
	Attempt   uint16 /* How many resends the client has done (count up from zero for same sequence number) */
	SeqNumber uint32 /* Client's sequence number */
	Pad       [396]byte
}

// TrackingPayload is the payload for tracking replies (`RPY_Tracking`)
type trackingPayload struct {
	RefID              uint32
	Ip                 ipAddr
	Stratum            uint16
	LeapStatus         uint16
	RefTime            chronyTimespec
	CurrentCorrection  chronyFloat
	LastOffset         chronyFloat
	RmsOffset          chronyFloat
	FreqPpm            chronyFloat
	ResidFreqPpm       chronyFloat
	SkewPpm            chronyFloat
	RootDelay          chronyFloat
	RootDispersion     chronyFloat
	LastUpdateInterval chronyFloat
}

const (
	IpaddrInet4 = uint16(1)
	IpaddrInet6 = uint16(2)
)

type ipAddr struct {
	IPAddrHigh uint64
	IPAddrLow  uint64
	Family     uint16
	Pad        uint16
}

func (tracking *trackingPayload) String() string {
	return fmt.Sprintf(
		"RefID: %d, ActivictServer: %s, Stratum: %d, RefTime: %s, CurrentCorrection: %f, "+
			"FreqPpm: %f, SkewPpm: %f, RootDelay: %f, "+
			"RootDispersion: %f, LeapStatus: %d, LastUpdateInterval: %f, "+
			"LastOffset: %f, CurrentCorrection: %f",
		tracking.RefID, tracking.Ip.String(), tracking.Stratum, tracking.RefTime.Time().Format(time.RFC3339),
		tracking.CurrentCorrection.Float64(), tracking.FreqPpm.Float64(), tracking.SkewPpm.Float64(),
		tracking.RootDelay.Float64(), tracking.RootDispersion.Float64(), tracking.LeapStatus,
		tracking.LastUpdateInterval.Float64(), tracking.LastOffset.Float64(), tracking.CurrentCorrection.Float64(),
	)
}

func (ia ipAddr) Ip() net.IP {
	if ia.Family == IpaddrInet4 {
		m := uint32(ia.IPAddrHigh >> (32))
		var ip [4]uint8
		for i := 0; i < 4; i++ {
			ip[i] = uint8(m % 0x100)
			m = m / 0x100
		}
		return net.IPv4(ip[3], ip[2], ip[1], ip[0])
	}

	if ia.Family == IpaddrInet6 {
		res := make(net.IP, net.IPv6len)
		h := ia.IPAddrHigh
		for i := 7; i >= 0; i-- {
			res[i] = byte(h % 0x100)
			h = h / 0x100
		}
		l := ia.IPAddrLow
		for i := 7; i >= 0; i-- {
			res[i+8] = byte(l % 0x100)
			l = l / 0x100
		}

		return res
	}

	return net.IPv4zero
}

func (ia ipAddr) String() string {
	return ia.Ip().String()
}

// ActivityPayload is the payload for activity replies (`RPY_Activity`)
type activityPayload struct {
	Online       int32
	Offline      int32
	BurstOnline  int32
	BurstOffline int32
	Unresolved   int32
}

func (activity *activityPayload) String() string {
	return fmt.Sprintf("Online: %d, Offline: %d, BurstOnline: %d, BurstOffline: %d, Unresolved: %d",
		activity.Online, activity.Offline, activity.BurstOnline, activity.BurstOffline, activity.Unresolved,
	)
}

// replyPacket is the common header for all replies
// chrony version 4.1
type replyPacket struct {
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

// chronyTimespec is the custom chrony timespec type (`Timespec`)
type chronyTimespec struct {
	TvSecHigh uint32
	TvSecLow  uint32
	TvNSec    uint32
}

func (ct chronyTimespec) Time() time.Time {
	var nsec = uint32(999999999)
	if ct.TvNSec < nsec {
		nsec = ct.TvNSec
	}

	return time.Unix(int64(uint64(ct.TvSecHigh)<<32+uint64(ct.TvSecLow)), int64(nsec))
}

// EpochSeconds returns the number of seconds since epoch
func (ct chronyTimespec) EpochSeconds() float64 {
	ts := uint64(ct.TvSecHigh) << 32
	ts += uint64(ct.TvSecLow)
	return float64(ts)
}

/* 32-bit floating-point format consisting of 7-bit signed exponent
   and 25-bit signed coefficient without hidden bit.
   The result is calculated as: 2^(exp - 25) * coef */
type chronyFloat int32

// Float64 does magic to decode float from int32.
// Code is copied and translated to Go from original C sources.
func (cf chronyFloat) Float64() float64 {
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
func (cf chronyFloat) Int64() int64 {
	return int64(cf.Float64() * scaleFactor)
}
