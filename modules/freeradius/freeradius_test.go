package freeradius

import (
	"context"
	"errors"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"layeh.com/radius"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestFreeradius_Init(t *testing.T) {
	mod := New()

	mod.Address = ""
	mod.Port = 0
	mod.Secret = ""

	assert.True(mod, mod.Init())
	assert.Equal(t, defaultAddress, mod.Address)
	assert.Equal(t, defaultPort, mod.Port)
	assert.Equal(t, defaultSecret, mod.Secret)
}

func TestFreeradius_Check(t *testing.T) {
	mod := New()

	mod.exchanger = &mockErrorRadiusServer{}
	assert.False(t, mod.Check())

	mod.exchanger = &mockOKRadiusServer{}
	assert.True(t, mod.Check())
}

func TestFreeradius_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestFreeradius_Collect(t *testing.T) {
	mod := New()

	mod.exchanger = &mockOKRadiusServer{}

	expected := map[string]int64{
		"access-requests":               10,
		"access-accepts":                10,
		"access-rejects":                10,
		"access-challenges":             10,
		"auth-responses":                10,
		"auth-duplicate-requests":       10,
		"auth-malformed-requests":       10,
		"auth-invalid-requests":         10,
		"auth-dropped-requests":         10,
		"auth-unknown-types":            10,
		"proxy-access-requests":         10,
		"proxy-access-accepts":          10,
		"proxy-access-rejects":          10,
		"proxy-access-challenges":       10,
		"proxy-auth-responses":          10,
		"proxy-auth-duplicate-requests": 10,
		"proxy-auth-malformed-requests": 10,
		"proxy-auth-invalid-requests":   10,
		"proxy-auth-dropped-requests":   10,
		"proxy-auth-unknown-types":      10,
		"accounting-requests":           10,
		"accounting-responses":          10,
		"acct-duplicate-requests":       10,
		"acct-malformed-requests":       10,
		"acct-invalid-requests":         10,
		"acct-dropped-requests":         10,
		"acct-unknown-types":            10,
		"proxy-accounting-requests":     10,
		"proxy-accounting-responses":    10,
		"proxy-acct-duplicate-requests": 10,
		"proxy-acct-malformed-requests": 10,
		"proxy-acct-invalid-requests":   10,
		"proxy-acct-dropped-requests":   10,
		"proxy-acct-unknown-types":      10,
	}

	assert.Equal(t, expected, mod.Collect())
}

func TestFreeradius_Cleanup(t *testing.T) {
	New().Cleanup()
}

type mockOKRadiusServer struct{}

func (m mockOKRadiusServer) Exchange(ctx context.Context, packet *radius.Packet, address string) (*radius.Packet, error) {
	resp := radius.New(radius.CodeAccessAccept, []byte(defaultSecret))
	resp.Code = radius.CodeAccessAccept

	_ = FreeRADIUSTotalAccessRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAccessAccepts_Add(resp, 10)
	_ = FreeRADIUSTotalAccessRejects_Add(resp, 10)
	_ = FreeRADIUSTotalAccessChallenges_Add(resp, 10)
	_ = FreeRADIUSTotalAuthResponses_Add(resp, 10)
	_ = FreeRADIUSTotalAuthDuplicateRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAuthMalformedRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAuthInvalidRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAuthDroppedRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAuthUnknownTypes_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAccessRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAccessAccepts_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAccessRejects_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAccessChallenges_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAuthResponses_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAuthDuplicateRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAuthMalformedRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAuthInvalidRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAuthDroppedRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAuthUnknownTypes_Add(resp, 10)
	_ = FreeRADIUSTotalAccountingRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAccountingResponses_Add(resp, 10)
	_ = FreeRADIUSTotalAcctDuplicateRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAcctMalformedRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAcctInvalidRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAcctDroppedRequests_Add(resp, 10)
	_ = FreeRADIUSTotalAcctUnknownTypes_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAccountingRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAccountingResponses_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAcctDuplicateRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAcctMalformedRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAcctInvalidRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAcctDroppedRequests_Add(resp, 10)
	_ = FreeRADIUSTotalProxyAcctUnknownTypes_Add(resp, 10)

	return resp, nil
}

type mockErrorRadiusServer struct{}

func (m mockErrorRadiusServer) Exchange(ctx context.Context, packet *radius.Packet, address string) (*radius.Packet, error) {
	return nil, errors.New("mock radius error")
}
