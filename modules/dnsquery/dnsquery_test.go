package dnsquery

import (
	"errors"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/modules"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*modules.Module)(nil), New())
}

func TestDNSQuery_Init(t *testing.T) {
	mod := New()

	// NG case
	assert.False(t, mod.Init())
	// OK case
	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	require.True(t, mod.Init())
	assert.NotZero(t, mod.servers)
	assert.NotZero(t, mod.workers)
}

func TestDNSQuery_Check(t *testing.T) {
	assert.True(t, New().Check())
}

func TestDNSQuery_Charts(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Charts())

	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	require.True(t, mod.Init())
	charts := mod.Charts()
	assert.True(t, charts.Get("query_time").HasDim("8_8_8_8"))
}

func TestDNSQuery_Cleanup(t *testing.T) {
	mod := New()

	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	require.True(t, mod.Init())
	assert.NotZero(t, mod.servers)
	assert.NotZero(t, mod.workers)
	mod.Cleanup()
	assert.Len(t, mod.workers, 0)
}

func TestDNSQuery_Collect(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	mod.exchangerFactory = func(network string, duration time.Duration) exchanger {
		return okMockExchanger{}
	}

	require.True(t, mod.Init())
	require.True(t, mod.Check())

	assert.Equal(
		t,
		map[string]int64{"8_8_8_8": 1000000000},
		mod.Collect(),
	)
}

func TestDNSQuery_Collect_Error(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	mod.exchangerFactory = func(network string, duration time.Duration) exchanger {
		return errMockExchanger{}
	}

	require.True(t, mod.Init())
	require.True(t, mod.Check())

	assert.Len(
		t,
		mod.Collect(),
		0,
	)
}

type okMockExchanger struct{}

func (m okMockExchanger) Exchange(msg *dns.Msg, address string) (response *dns.Msg, rtt time.Duration, err error) {
	return nil, time.Second, nil
}

type errMockExchanger struct{}

func (m errMockExchanger) Exchange(msg *dns.Msg, address string) (response *dns.Msg, rtt time.Duration, err error) {
	return nil, time.Second, errors.New("mock error")
}
