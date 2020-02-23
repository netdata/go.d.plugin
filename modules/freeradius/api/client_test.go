package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"layeh.com/radius"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New(Config{}))
}

func TestClient_Status(t *testing.T) {
	var c Client
	c.radiusClient = newOKMockFreeRADIUSClient()

	expected := Status{
		AccessRequests:             1,
		AccessAccepts:              2,
		AccessRejects:              3,
		AccessChallenges:           4,
		AuthResponses:              5,
		AuthDuplicateRequests:      6,
		AuthMalformedRequests:      7,
		AuthInvalidRequests:        8,
		AuthDroppedRequests:        9,
		AuthUnknownTypes:           10,
		AccountingRequests:         11,
		AccountingResponses:        12,
		AcctDuplicateRequests:      13,
		AcctMalformedRequests:      14,
		AcctInvalidRequests:        15,
		AcctDroppedRequests:        16,
		AcctUnknownTypes:           17,
		ProxyAccessRequests:        18,
		ProxyAccessAccepts:         19,
		ProxyAccessRejects:         20,
		ProxyAccessChallenges:      21,
		ProxyAuthResponses:         22,
		ProxyAuthDuplicateRequests: 23,
		ProxyAuthMalformedRequests: 24,
		ProxyAuthInvalidRequests:   25,
		ProxyAuthDroppedRequests:   26,
		ProxyAuthUnknownTypes:      27,
		ProxyAccountingRequests:    28,
		ProxyAccountingResponses:   29,
		ProxyAcctDuplicateRequests: 30,
		ProxyAcctMalformedRequests: 31,
		ProxyAcctInvalidRequests:   32,
		ProxyAcctDroppedRequests:   33,
		ProxyAcctUnknownTypes:      34,
	}

	s, err := c.Status()

	require.NoError(t, err)
	assert.Equal(t, expected, *s)
}

func TestClient_Status_ReturnsErrorIfClientExchangeReturnsError(t *testing.T) {
	var c Client
	c.radiusClient = newErrorMockFreeRADIUSClient()

	s, err := c.Status()

	assert.Nil(t, s)
	assert.Error(t, err)
}

func TestClient_Status_ReturnsErrorIfServerResponseHasBadStatus(t *testing.T) {
	var c Client
	c.radiusClient = newBadRespCodeMockFreeRADIUSClient()

	s, err := c.Status()

	assert.Nil(t, s)
	assert.Error(t, err)
}

type mockFreeRADIUSClient struct {
	errOnExchange bool
	badRespCode   bool
}

func newOKMockFreeRADIUSClient() *mockFreeRADIUSClient {
	return &mockFreeRADIUSClient{}
}

func newErrorMockFreeRADIUSClient() *mockFreeRADIUSClient {
	return &mockFreeRADIUSClient{errOnExchange: true}
}

func newBadRespCodeMockFreeRADIUSClient() *mockFreeRADIUSClient {
	return &mockFreeRADIUSClient{badRespCode: true}
}

func (m mockFreeRADIUSClient) Exchange(_ context.Context, _ *radius.Packet, _ string) (*radius.Packet, error) {
	if m.errOnExchange {
		return nil, errors.New("mock Exchange error")
	}
	resp := radius.New(radius.CodeAccessAccept, []byte("secret"))
	if m.badRespCode {
		resp.Code = radius.CodeAccessReject
	} else {
		resp.Code = radius.CodeAccessAccept
	}
	addValues(resp)
	return resp, nil
}

func addValues(resp *radius.Packet) {
	addAccessRequests(resp, 1)
	addAccessAccepts(resp, 2)
	addAccessRejects(resp, 3)
	addAccessChallenges(resp, 4)
	addAuthResponses(resp, 5)
	addAuthDuplicateRequests(resp, 6)
	addAuthMalformedRequests(resp, 7)
	addAuthInvalidRequests(resp, 8)
	addAuthDroppedRequests(resp, 9)
	addAuthUnknownTypes(resp, 10)
	addAccountingRequests(resp, 11)
	addAccountingResponses(resp, 12)
	addAcctDuplicateRequests(resp, 13)
	addAcctMalformedRequests(resp, 14)
	addAcctInvalidRequests(resp, 15)
	addAcctDroppedRequests(resp, 16)
	addAcctUnknownTypes(resp, 17)
	addProxyAccessRequests(resp, 18)
	addProxyAccessAccepts(resp, 19)
	addProxyAccessRejects(resp, 20)
	addProxyAccessChallenges(resp, 21)
	addProxyAuthResponses(resp, 22)
	addProxyAuthDuplicateRequests(resp, 23)
	addProxyAuthMalformedRequests(resp, 24)
	addProxyAuthInvalidRequests(resp, 25)
	addProxyAuthDroppedRequests(resp, 26)
	addProxyAuthUnknownTypes(resp, 27)
	addProxyAccountingRequests(resp, 28)
	addProxyAccountingResponses(resp, 29)
	addProxyAcctDuplicateRequests(resp, 30)
	addProxyAcctMalformedRequests(resp, 31)
	addProxyAcctInvalidRequests(resp, 32)
	addProxyAcctDroppedRequests(resp, 33)
	addProxyAcctUnknownTypes(resp, 34)
}

func addAccessRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAccessRequests_Add(resp, FreeRADIUSTotalAccessRequests(value))
}

func addAccessAccepts(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAccessAccepts_Add(resp, FreeRADIUSTotalAccessAccepts(value))
}

func addAccessRejects(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAccessRejects_Add(resp, FreeRADIUSTotalAccessRejects(value))
}

func addAccessChallenges(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAccessChallenges_Add(resp, FreeRADIUSTotalAccessChallenges(value))
}

func addAuthResponses(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAuthResponses_Add(resp, FreeRADIUSTotalAuthResponses(value))
}

func addAuthDuplicateRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAuthDuplicateRequests_Add(resp, FreeRADIUSTotalAuthDuplicateRequests(value))
}

func addAuthMalformedRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAuthMalformedRequests_Add(resp, FreeRADIUSTotalAuthMalformedRequests(value))
}

func addAuthInvalidRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAuthInvalidRequests_Add(resp, FreeRADIUSTotalAuthInvalidRequests(value))
}

func addAuthDroppedRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAuthDroppedRequests_Add(resp, FreeRADIUSTotalAuthDroppedRequests(value))
}

func addAuthUnknownTypes(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAuthUnknownTypes_Add(resp, FreeRADIUSTotalAuthUnknownTypes(value))
}

func addAccountingRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAccountingRequests_Add(resp, FreeRADIUSTotalAccountingRequests(value))
}

func addAccountingResponses(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAccountingResponses_Add(resp, FreeRADIUSTotalAccountingResponses(value))
}

func addAcctDuplicateRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAcctDuplicateRequests_Add(resp, FreeRADIUSTotalAcctDuplicateRequests(value))
}

func addAcctMalformedRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAcctMalformedRequests_Add(resp, FreeRADIUSTotalAcctMalformedRequests(value))
}

func addAcctInvalidRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAcctInvalidRequests_Add(resp, FreeRADIUSTotalAcctInvalidRequests(value))
}

func addAcctDroppedRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAcctDroppedRequests_Add(resp, FreeRADIUSTotalAcctDroppedRequests(value))
}

func addAcctUnknownTypes(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalAcctUnknownTypes_Add(resp, FreeRADIUSTotalAcctUnknownTypes(value))
}

func addProxyAccessRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAccessRequests_Add(resp, FreeRADIUSTotalProxyAccessRequests(value))
}

func addProxyAccessAccepts(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAccessAccepts_Add(resp, FreeRADIUSTotalProxyAccessAccepts(value))
}

func addProxyAccessRejects(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAccessRejects_Add(resp, FreeRADIUSTotalProxyAccessRejects(value))
}

func addProxyAccessChallenges(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAccessChallenges_Add(resp, FreeRADIUSTotalProxyAccessChallenges(value))
}

func addProxyAuthResponses(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAuthResponses_Add(resp, FreeRADIUSTotalProxyAuthResponses(value))
}

func addProxyAuthDuplicateRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAuthDuplicateRequests_Add(resp, FreeRADIUSTotalProxyAuthDuplicateRequests(value))
}

func addProxyAuthMalformedRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAuthMalformedRequests_Add(resp, FreeRADIUSTotalProxyAuthMalformedRequests(value))
}

func addProxyAuthInvalidRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAuthInvalidRequests_Add(resp, FreeRADIUSTotalProxyAuthInvalidRequests(value))
}

func addProxyAuthDroppedRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAuthDroppedRequests_Add(resp, FreeRADIUSTotalProxyAuthDroppedRequests(value))
}

func addProxyAuthUnknownTypes(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAuthUnknownTypes_Add(resp, FreeRADIUSTotalProxyAuthUnknownTypes(value))
}

func addProxyAccountingRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAccountingRequests_Add(resp, FreeRADIUSTotalProxyAccountingRequests(value))
}

func addProxyAccountingResponses(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAccountingResponses_Add(resp, FreeRADIUSTotalProxyAccountingResponses(value))
}

func addProxyAcctDuplicateRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAcctDuplicateRequests_Add(resp, FreeRADIUSTotalProxyAcctDuplicateRequests(value))
}

func addProxyAcctMalformedRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAcctMalformedRequests_Add(resp, FreeRADIUSTotalProxyAcctMalformedRequests(value))
}

func addProxyAcctInvalidRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAcctInvalidRequests_Add(resp, FreeRADIUSTotalProxyAcctInvalidRequests(value))
}

func addProxyAcctDroppedRequests(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAcctDroppedRequests_Add(resp, FreeRADIUSTotalProxyAcctDroppedRequests(value))
}

func addProxyAcctUnknownTypes(resp *radius.Packet, value int64) {
	_ = FreeRADIUSTotalProxyAcctUnknownTypes_Add(resp, FreeRADIUSTotalProxyAcctUnknownTypes(value))
}
