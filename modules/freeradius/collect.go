package freeradius

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"net"
	"strconv"

	"layeh.com/radius"
	"layeh.com/radius/rfc2869"
)

func (f *FreeRADIUS) collect() (map[string]int64, error) {
	packet, err := newStatusServerPacket(f.Secret)
	if err != nil {
		return nil, fmt.Errorf("error on creating StatusServer packet: %v", err)
	}

	resp, err := f.query(packet)
	if err != nil {
		return nil, err
	}

	return decodeResponse(resp), nil
}

func (f *FreeRADIUS) query(packet *radius.Packet) (*radius.Packet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.Timeout.Duration)
	defer cancel()

	address := net.JoinHostPort(f.Address, strconv.Itoa(f.Port))
	resp, err := f.client.Exchange(ctx, packet, address)
	if err != nil {
		return nil, fmt.Errorf("error on request to '%s': %s", address, err)
	}

	if resp.Code != radius.CodeAccessAccept {
		return nil, fmt.Errorf("'%s' returned response code %d", address, resp.Code)
	}
	return resp, nil
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
	if _, err := hash.Write(encode); err != nil {
		return nil, err
	}
	if err := rfc2869.MessageAuthenticator_Set(packet, hash.Sum(nil)); err != nil {
		return nil, err
	}
	return packet, nil
}

func decodeResponse(resp *radius.Packet) map[string]int64 {
	return map[string]int64{
		// authentication
		"access-requests":         int64(FreeRADIUSTotalAccessRequests_Get(resp)),
		"auth-responses":          int64(FreeRADIUSTotalAuthResponses_Get(resp)),
		"access-accepts":          int64(FreeRADIUSTotalAccessAccepts_Get(resp)),
		"access-rejects":          int64(FreeRADIUSTotalAccessRejects_Get(resp)),
		"access-challenges":       int64(FreeRADIUSTotalAccessChallenges_Get(resp)),
		"auth-duplicate-requests": int64(FreeRADIUSTotalAuthDuplicateRequests_Get(resp)),
		"auth-malformed-requests": int64(FreeRADIUSTotalAuthMalformedRequests_Get(resp)),
		"auth-invalid-requests":   int64(FreeRADIUSTotalAuthInvalidRequests_Get(resp)),
		"auth-dropped-requests":   int64(FreeRADIUSTotalAuthDroppedRequests_Get(resp)),
		"auth-unknown-types":      int64(FreeRADIUSTotalAuthUnknownTypes_Get(resp)),

		// accounting
		"accounting-requests":     int64(FreeRADIUSTotalAccountingRequests_Get(resp)),
		"accounting-responses":    int64(FreeRADIUSTotalAccountingResponses_Get(resp)),
		"acct-duplicate-requests": int64(FreeRADIUSTotalAcctDuplicateRequests_Get(resp)),
		"acct-malformed-requests": int64(FreeRADIUSTotalAcctMalformedRequests_Get(resp)),
		"acct-invalid-requests":   int64(FreeRADIUSTotalAcctInvalidRequests_Get(resp)),
		"acct-dropped-requests":   int64(FreeRADIUSTotalAcctDroppedRequests_Get(resp)),
		"acct-unknown-types":      int64(FreeRADIUSTotalAcctUnknownTypes_Get(resp)),

		// proxy authentication
		"proxy-access-requests":         int64(FreeRADIUSTotalProxyAccessRequests_Get(resp)),
		"proxy-auth-responses":          int64(FreeRADIUSTotalProxyAuthResponses_Get(resp)),
		"proxy-access-accepts":          int64(FreeRADIUSTotalProxyAccessAccepts_Get(resp)),
		"proxy-access-rejects":          int64(FreeRADIUSTotalProxyAccessRejects_Get(resp)),
		"proxy-access-challenges":       int64(FreeRADIUSTotalProxyAccessChallenges_Get(resp)),
		"proxy-auth-duplicate-requests": int64(FreeRADIUSTotalProxyAuthDuplicateRequests_Get(resp)),
		"proxy-auth-malformed-requests": int64(FreeRADIUSTotalProxyAuthMalformedRequests_Get(resp)),
		"proxy-auth-invalid-requests":   int64(FreeRADIUSTotalProxyAuthInvalidRequests_Get(resp)),
		"proxy-auth-dropped-requests":   int64(FreeRADIUSTotalProxyAuthDroppedRequests_Get(resp)),
		"proxy-auth-unknown-types":      int64(FreeRADIUSTotalProxyAuthUnknownTypes_Get(resp)),

		// proxy accounting
		"proxy-accounting-requests":     int64(FreeRADIUSTotalProxyAccountingRequests_Get(resp)),
		"proxy-accounting-responses":    int64(FreeRADIUSTotalProxyAccountingResponses_Get(resp)),
		"proxy-acct-duplicate-requests": int64(FreeRADIUSTotalProxyAcctDuplicateRequests_Get(resp)),
		"proxy-acct-malformed-requests": int64(FreeRADIUSTotalProxyAcctMalformedRequests_Get(resp)),
		"proxy-acct-invalid-requests":   int64(FreeRADIUSTotalProxyAcctInvalidRequests_Get(resp)),
		"proxy-acct-dropped-requests":   int64(FreeRADIUSTotalProxyAcctDroppedRequests_Get(resp)),
		"proxy-acct-unknown-types":      int64(FreeRADIUSTotalProxyAcctUnknownTypes_Get(resp)),
	}
}
