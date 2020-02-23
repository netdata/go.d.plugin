package freeradius

import (
	"context"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"layeh.com/radius"
)

var (
	allStatisticsData, _ = ioutil.ReadFile("testdata/all.txt")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, allStatisticsData)
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestFreeRADIUS_Init(t *testing.T) {
	freeRADIUS := New()

	assert.True(t, freeRADIUS.Init())
}

func TestFreeRADIUS_Init_ReturnsFalseIfAddressNotSet(t *testing.T) {
	freeRADIUS := New()
	freeRADIUS.Address = ""

	assert.False(t, freeRADIUS.Init())
}

func TestFreeRADIUS_Init_ReturnsFalseIfPortNotSet(t *testing.T) {
	freeRADIUS := New()
	freeRADIUS.Port = 0

	assert.False(t, freeRADIUS.Init())
}

func TestFreeRADIUS_Init_ReturnsFalseIfSecretNotSet(t *testing.T) {
	freeRADIUS := New()
	freeRADIUS.Secret = ""

	assert.False(t, freeRADIUS.Init())
}

func TestFreeRADIUS_Check(t *testing.T) {
	freeRADIUS := New()
	freeRADIUS.client = newOKMockFreeRADIUSClient()

	assert.True(t, freeRADIUS.Check())
}

func TestFreeRADIUS_Check_ReturnsFalseIfClientExchangeReturnsError(t *testing.T) {
	freeRADIUS := New()
	freeRADIUS.client = newErrorMockFreeRADIUSClient()

	assert.False(t, freeRADIUS.Check())
}

func TestFreeRADIUS_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestFreeRADIUS_Collect(t *testing.T) {
	freeRADIUS := New()
	freeRADIUS.client = newOKMockFreeRADIUSClient()

	expected := map[string]int64{
		"access-requests":               1,
		"access-accepts":                2,
		"access-rejects":                3,
		"access-challenges":             4,
		"auth-responses":                5,
		"auth-duplicate-requests":       6,
		"auth-malformed-requests":       7,
		"auth-invalid-requests":         8,
		"auth-dropped-requests":         9,
		"auth-unknown-types":            10,
		"accounting-requests":           11,
		"accounting-responses":          12,
		"acct-duplicate-requests":       13,
		"acct-malformed-requests":       14,
		"acct-invalid-requests":         15,
		"acct-dropped-requests":         16,
		"acct-unknown-types":            17,
		"proxy-access-requests":         18,
		"proxy-access-accepts":          19,
		"proxy-access-rejects":          20,
		"proxy-access-challenges":       21,
		"proxy-auth-duplicate-requests": 22,
		"proxy-auth-responses":          22,
		"proxy-auth-malformed-requests": 23,
		"proxy-auth-invalid-requests":   24,
		"proxy-auth-dropped-requests":   25,
		"proxy-auth-unknown-types":      26,
		"proxy-accounting-requests":     27,
		"proxy-accounting-responses":    28,
		"proxy-acct-duplicate-requests": 29,
		"proxy-acct-malformed-requests": 30,
		"proxy-acct-invalid-requests":   31,
		"proxy-acct-dropped-requests":   32,
		"proxy-acct-unknown-types":      33,
	}

	assert.Equal(t, expected, freeRADIUS.Collect())
}

func TestFreeRADIUS_Collect_ReturnsNilIfClientExchangeReturnsError(t *testing.T) {
	freeRADIUS := New()
	freeRADIUS.client = newErrorMockFreeRADIUSClient()

	assert.Nil(t, freeRADIUS.Collect())
}

func TestFreeRADIUS_Cleanup(t *testing.T) {
	New().Cleanup()
}

func newOKMockFreeRADIUSClient() *mockFreeRADIUSClient {
	return &mockFreeRADIUSClient{
		errOnExchange: false,
		secret:        "adminsercet",
		response:      allStatisticsData,
	}
}

func newErrorMockFreeRADIUSClient() *mockFreeRADIUSClient {
	return &mockFreeRADIUSClient{
		errOnExchange: true,
	}
}

type mockFreeRADIUSClient struct {
	errOnExchange bool
	secret        string
	response      []byte
}

func (m mockFreeRADIUSClient) Exchange(_ context.Context, _ *radius.Packet, _ string) (*radius.Packet, error) {
	if m.errOnExchange {
		return nil, errors.New("mock Exchange error")
	}
	resp := radius.New(radius.CodeAccessAccept, []byte(m.secret))
	resp.Code = radius.CodeAccessAccept
	m.setResponse(resp)
	return resp, nil
}

func (m mockFreeRADIUSClient) setResponse(resp *radius.Packet) {
	for _, line := range strings.Split(string(m.response), "\n") {
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		value, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			continue
		}

		name := strings.TrimSpace(parts[0])
		switch name {
		case "FreeRADIUS-Total-Access-Requests":
			_ = FreeRADIUSTotalAccessRequests_Add(resp, FreeRADIUSTotalAccessRequests(value))
		case "FreeRADIUS-Total-Access-Accepts":
			_ = FreeRADIUSTotalAccessAccepts_Add(resp, FreeRADIUSTotalAccessAccepts(value))
		case "FreeRADIUS-Total-Access-Rejects":
			_ = FreeRADIUSTotalAccessRejects_Add(resp, FreeRADIUSTotalAccessRejects(value))
		case "FreeRADIUS-Total-Access-Challenges":
			_ = FreeRADIUSTotalAccessChallenges_Add(resp, FreeRADIUSTotalAccessChallenges(value))
		case "FreeRADIUS-Total-Auth-Responses":
			_ = FreeRADIUSTotalAuthResponses_Add(resp, FreeRADIUSTotalAuthResponses(value))
		case "FreeRADIUS-Total-Auth-Duplicate-Requests":
			_ = FreeRADIUSTotalAuthDuplicateRequests_Add(resp, FreeRADIUSTotalAuthDuplicateRequests(value))
		case "FreeRADIUS-Total-Auth-Malformed-Requests":
			_ = FreeRADIUSTotalAuthMalformedRequests_Add(resp, FreeRADIUSTotalAuthMalformedRequests(value))
		case "FreeRADIUS-Total-Auth-Invalid-Requests":
			_ = FreeRADIUSTotalAuthInvalidRequests_Add(resp, FreeRADIUSTotalAuthInvalidRequests(value))
		case "FreeRADIUS-Total-Auth-Dropped-Requests":
			_ = FreeRADIUSTotalAuthDroppedRequests_Add(resp, FreeRADIUSTotalAuthDroppedRequests(value))
		case "FreeRADIUS-Total-Auth-Unknown-Types":
			_ = FreeRADIUSTotalAuthUnknownTypes_Add(resp, FreeRADIUSTotalAuthUnknownTypes(value))

		case "FreeRADIUS-Total-Accounting-Requests":
			_ = FreeRADIUSTotalAccountingRequests_Add(resp, FreeRADIUSTotalAccountingRequests(value))
		case "FreeRADIUS-Total-Accounting-Responses":
			_ = FreeRADIUSTotalAccountingResponses_Add(resp, FreeRADIUSTotalAccountingResponses(value))
		case "FreeRADIUS-Total-Acct-Duplicate-Requests":
			_ = FreeRADIUSTotalAcctDuplicateRequests_Add(resp, FreeRADIUSTotalAcctDuplicateRequests(value))
		case "FreeRADIUS-Total-Acct-Malformed-Requests":
			_ = FreeRADIUSTotalAcctMalformedRequests_Add(resp, FreeRADIUSTotalAcctMalformedRequests(value))
		case "FreeRADIUS-Total-Acct-Invalid-Requests":
			_ = FreeRADIUSTotalAcctInvalidRequests_Add(resp, FreeRADIUSTotalAcctInvalidRequests(value))
		case "FreeRADIUS-Total-Acct-Dropped-Requests":
			_ = FreeRADIUSTotalAcctDroppedRequests_Add(resp, FreeRADIUSTotalAcctDroppedRequests(value))
		case "FreeRADIUS-Total-Acct-Unknown-Types":
			_ = FreeRADIUSTotalAcctUnknownTypes_Add(resp, FreeRADIUSTotalAcctUnknownTypes(value))

		case "FreeRADIUS-Total-Proxy-Access-Requests":
			_ = FreeRADIUSTotalProxyAccessRequests_Add(resp, FreeRADIUSTotalProxyAccessRequests(value))
		case "FreeRADIUS-Total-Proxy-Access-Accepts":
			_ = FreeRADIUSTotalProxyAccessAccepts_Add(resp, FreeRADIUSTotalProxyAccessAccepts(value))
		case "FreeRADIUS-Total-Proxy-Access-Rejects":
			_ = FreeRADIUSTotalProxyAccessRejects_Add(resp, FreeRADIUSTotalProxyAccessRejects(value))
		case "FreeRADIUS-Total-Proxy-Access-Challenges":
			_ = FreeRADIUSTotalProxyAccessChallenges_Add(resp, FreeRADIUSTotalProxyAccessChallenges(value))
		case "FreeRADIUS-Total-Proxy-Auth-Responses":
			_ = FreeRADIUSTotalProxyAuthResponses_Add(resp, FreeRADIUSTotalProxyAuthResponses(value))
		case "FreeRADIUS-Total-Proxy-Auth-Duplicate-Requests":
			_ = FreeRADIUSTotalProxyAuthDuplicateRequests_Add(resp, FreeRADIUSTotalProxyAuthDuplicateRequests(value))
		case "FreeRADIUS-Total-Proxy-Auth-Malformed-Requests":
			_ = FreeRADIUSTotalProxyAuthMalformedRequests_Add(resp, FreeRADIUSTotalProxyAuthMalformedRequests(value))
		case "FreeRADIUS-Total-Proxy-Auth-Invalid-Requests":
			_ = FreeRADIUSTotalProxyAuthInvalidRequests_Add(resp, FreeRADIUSTotalProxyAuthInvalidRequests(value))
		case "FreeRADIUS-Total-Proxy-Auth-Dropped-Requests":
			_ = FreeRADIUSTotalProxyAuthDroppedRequests_Add(resp, FreeRADIUSTotalProxyAuthDroppedRequests(value))
		case "FreeRADIUS-Total-Proxy-Auth-Unknown-Types":
			_ = FreeRADIUSTotalProxyAuthUnknownTypes_Add(resp, FreeRADIUSTotalProxyAuthUnknownTypes(value))

		case "FreeRADIUS-Total-Proxy-Accounting-Requests":
			_ = FreeRADIUSTotalProxyAccountingRequests_Add(resp, FreeRADIUSTotalProxyAccountingRequests(value))
		case "FreeRADIUS-Total-Proxy-Accounting-Responses":
			_ = FreeRADIUSTotalProxyAccountingResponses_Add(resp, FreeRADIUSTotalProxyAccountingResponses(value))
		case "FreeRADIUS-Total-Proxy-Acct-Duplicate-Requests":
			_ = FreeRADIUSTotalProxyAcctDuplicateRequests_Add(resp, FreeRADIUSTotalProxyAcctDuplicateRequests(value))
		case "FreeRADIUS-Total-Proxy-Acct-Malformed-Requests":
			_ = FreeRADIUSTotalProxyAcctMalformedRequests_Add(resp, FreeRADIUSTotalProxyAcctMalformedRequests(value))
		case "FreeRADIUS-Total-Proxy-Acct-Invalid-Requests":
			_ = FreeRADIUSTotalProxyAcctInvalidRequests_Add(resp, FreeRADIUSTotalProxyAcctInvalidRequests(value))
		case "FreeRADIUS-Total-Proxy-Acct-Dropped-Requests":
			_ = FreeRADIUSTotalProxyAcctDroppedRequests_Add(resp, FreeRADIUSTotalProxyAcctDroppedRequests(value))
		case "FreeRADIUS-Total-Proxy-Acct-Unknown-Types":
			_ = FreeRADIUSTotalProxyAcctUnknownTypes_Add(resp, FreeRADIUSTotalProxyAcctUnknownTypes(value))
		}
	}
}
