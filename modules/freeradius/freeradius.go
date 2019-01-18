package freeradius

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"net"
	"strconv"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"

	"layeh.com/radius"
	"layeh.com/radius/rfc2869"
)

func init() {
	creator := modules.Creator{
		DisabledByDefault: true,
		Create:            func() modules.Module { return New() },
	}

	modules.Register("freeradius", creator)
}

const (
	defaultAddress = "127.0.0.1"
	defaultPort    = 18121
	defaultSecret  = "adminsecret"
)

// New creates Freeradius with default values
func New() *Freeradius {
	return &Freeradius{
		Address: defaultAddress,
		Port:    defaultPort,
		Secret:  defaultSecret,
		Timeout: web.Duration{Duration: time.Second},

		exchanger: &radius.Client{
			Retry:           time.Second,
			MaxPacketErrors: 10,
		},
	}
}

type exchanger interface {
	Exchange(ctx context.Context, packet *radius.Packet, address string) (*radius.Packet, error)
}

// Freeradius freeradius module
type Freeradius struct {
	modules.Base

	Address string
	Port    int
	Secret  string
	Timeout web.Duration

	exchanger exchanger
}

// Cleanup makes cleanup
func (Freeradius) Cleanup() {}

// Init makes initialization
func (f *Freeradius) Init() bool {
	if f.Address == "" {
		f.Address = defaultAddress
	}
	if f.Port == 0 {
		f.Port = defaultPort
	}
	if f.Secret == "" {
		f.Secret = defaultSecret
	}

	return true
}

// Check makes check
func (f Freeradius) Check() bool {
	return len(f.Collect()) > 0
}

// Charts creates Charts
func (Freeradius) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (f *Freeradius) Collect() map[string]int64 {
	packet, err := newStatusServerPacket(f.Secret)

	if err != nil {
		f.Errorf("error on creating StatusServer packet : %s", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), f.Timeout.Duration)
	defer cancel()

	resp, err := f.exchanger.Exchange(ctx, packet, net.JoinHostPort(f.Address, strconv.Itoa(f.Port)))

	if err != nil {
		f.Errorf("error on request to %s : %s", net.JoinHostPort(f.Address, strconv.Itoa(f.Port)), err)
		return nil
	}

	if resp.Code != radius.CodeAccessAccept {
		f.Errorf("%s returned response code %d", net.JoinHostPort(f.Address, strconv.Itoa(f.Port)), resp.Code)
		return nil
	}

	return decodeServerResponse(resp)
}

func decodeServerResponse(resp *radius.Packet) map[string]int64 {
	return map[string]int64{
		"access-requests":               int64(FreeRADIUSTotalAccessRequests_Get(resp)),
		"access-accepts":                int64(FreeRADIUSTotalAccessAccepts_Get(resp)),
		"access-rejects":                int64(FreeRADIUSTotalAccessRequests_Get(resp)),
		"access-challenges":             int64(FreeRADIUSTotalAccessChallenges_Get(resp)),
		"auth-responses":                int64(FreeRADIUSTotalAuthResponses_Get(resp)),
		"auth-duplicate-requests":       int64(FreeRADIUSTotalAuthDuplicateRequests_Get(resp)),
		"auth-malformed-requests":       int64(FreeRADIUSTotalAuthMalformedRequests_Get(resp)),
		"auth-invalid-requests":         int64(FreeRADIUSTotalAuthInvalidRequests_Get(resp)),
		"auth-dropped-requests":         int64(FreeRADIUSTotalAuthDroppedRequests_Get(resp)),
		"auth-unknown-types":            int64(FreeRADIUSTotalAuthUnknownTypes_Get(resp)),
		"proxy-access-requests":         int64(FreeRADIUSTotalProxyAccessRequests_Get(resp)),
		"proxy-access-accepts":          int64(FreeRADIUSTotalProxyAccessAccepts_Get(resp)),
		"proxy-access-rejects":          int64(FreeRADIUSTotalProxyAccessRequests_Get(resp)),
		"proxy-access-challenges":       int64(FreeRADIUSTotalProxyAccessChallenges_Get(resp)),
		"proxy-auth-responses":          int64(FreeRADIUSTotalProxyAuthResponses_Get(resp)),
		"proxy-auth-duplicate-requests": int64(FreeRADIUSTotalProxyAuthDuplicateRequests_Get(resp)),
		"proxy-auth-malformed-requests": int64(FreeRADIUSTotalProxyAuthMalformedRequests_Get(resp)),
		"proxy-auth-invalid-requests":   int64(FreeRADIUSTotalProxyAuthInvalidRequests_Get(resp)),
		"proxy-auth-dropped-requests":   int64(FreeRADIUSTotalProxyAuthDroppedRequests_Get(resp)),
		"proxy-auth-unknown-types":      int64(FreeRADIUSTotalProxyAuthUnknownTypes_Get(resp)),
		"accounting-requests":           int64(FreeRADIUSTotalAccountingRequests_Get(resp)),
		"accounting-responses":          int64(FreeRADIUSTotalAccountingResponses_Get(resp)),
		"acct-duplicate-requests":       int64(FreeRADIUSTotalAcctDuplicateRequests_Get(resp)),
		"acct-malformed-requests":       int64(FreeRADIUSTotalAcctMalformedRequests_Get(resp)),
		"acct-invalid-requests":         int64(FreeRADIUSTotalAcctInvalidRequests_Get(resp)),
		"acct-dropped-requests":         int64(FreeRADIUSTotalAcctDroppedRequests_Get(resp)),
		"acct-unknown-types":            int64(FreeRADIUSTotalAcctUnknownTypes_Get(resp)),
		"proxy-accounting-requests":     int64(FreeRADIUSTotalProxyAccountingRequests_Get(resp)),
		"proxy-accounting-responses":    int64(FreeRADIUSTotalProxyAccountingResponses_Get(resp)),
		"proxy-acct-duplicate-requests": int64(FreeRADIUSTotalProxyAcctDuplicateRequests_Get(resp)),
		"proxy-acct-malformed-requests": int64(FreeRADIUSTotalProxyAcctMalformedRequests_Get(resp)),
		"proxy-acct-invalid-requests":   int64(FreeRADIUSTotalProxyAcctInvalidRequests_Get(resp)),
		"proxy-acct-dropped-requests":   int64(FreeRADIUSTotalProxyAcctDroppedRequests_Get(resp)),
		"proxy-acct-unknown-types":      int64(FreeRADIUSTotalProxyAcctUnknownTypes_Get(resp)),
	}
}

func newStatusServerPacket(secret string) (*radius.Packet, error) {
	// https://wiki.freeradius.org/config/Status#status-of-freeradius-server
	packet := radius.New(radius.CodeStatusServer, []byte(secret))

	if err := FreeRADIUSStatisticsType_Set(packet, FreeRADIUSStatisticsType_Value_All); err != nil {
		return nil, err
	}

	if err := rfc2869.MessageAuthenticator_Set(packet, make([]byte, 16)); err != nil {
		return nil, err
	}

	hash := hmac.New(md5.New, packet.Secret)

	encode, err := packet.Encode()

	if err != nil {
		return nil, err
	}

	hash.Write(encode)

	if err := rfc2869.MessageAuthenticator_Set(packet, hash.Sum(nil)); err != nil {
		return nil, err
	}

	return packet, nil
}
